package mongo

import (
	"log"
	"time"
	"github.com/pkg/errors"
	"github.com/gocql/gocql"
	"ms/transmission/core"
)

const (
	CREATE_KEYSPACE    = " CREATE KEYSPACE IF NOT EXISTS " + KEYSPACE_NAME + " WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };"
	CREATE_TABLE       = "create table if not exists transmission.clients (idMeeting text, idUser int, idSession int, tokenData text, media blob, primary key (idMeeting, idUser));"
	KEYSPACE_NAME      = "transmission"
	CASSANDRA_URL      = "CASSANDRA_URL"
	CASSANDRA_USERNAME = "CASSANDRA_USERNAME"
	CASSANDRA_PASSWORD = "CASSANDRA_PASSWORD"
	PAGE_SIZE          = 10
)

type CassandraRepository struct {
	session  *gocql.Session
}

func createCluster(host string, keyspace string, authentication gocql.PasswordAuthenticator) *gocql.ClusterConfig {
	cluster := gocql.NewCluster(host)
	cluster.Authenticator = authentication
	createKeyspace(keyspace, cluster)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.One
	cluster.Timeout = time.Second * 10
	return cluster
}

func createKeyspace(keyspace string, cluster *gocql.ClusterConfig) {
	session, err := cluster.CreateSession()
	defer session.Close()
	if err != nil {
		log.Fatal("FATAL", err)
	}
	if err := session.Query(CREATE_KEYSPACE).Exec(); err != nil {
		log.Fatal("FATAL", err)
	}
	log.Println("INFO", "Configurado keyspace: "+keyspace)
}

func createSessionTable(session *gocql.Session) {
	log.Println("INFO", "Creando tabla si no existe...")
	err := session.Query(CREATE_TABLE).Exec(); 
	if err != nil {
		log.Println("ERROR", "Error creando tabla", err)
	}
}

func NewRepository(url, dbName, user, pass string) (core.ClientRepository, error) {
	repo := &CassandraRepository{}
	auth := gocql.PasswordAuthenticator{Username: user, Password: pass}
	cluster := createCluster(url, KEYSPACE_NAME, auth)
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRepository")
	}
	createSessionTable(session)
	session.SetPageSize(PAGE_SIZE)
	repo.session = session
	return repo, nil
}

func (repository *CassandraRepository) List(token string) ([]*core.Client, error) {
	var clientList []*core.Client
	var meetingKey string

	var meetingQuery = "select idMeeting from transmission.clients where tokenData=? ALLOW FILTERING"
	meetingOutput := repository.session.Query(meetingQuery, token).Iter()
	meetingOutput.Scan(&meetingKey)
	
	var selectQuery = "select * from transmission.clients where idMeeting=? ALLOW FILTERING"
	m := map[string]interface{}{}
	query := repository.session.Query(selectQuery, &meetingKey).Iter()
	if query == nil {
		log.Println("ERROR", "Error obteniendo clientes")
		return nil, errors.Wrap(nil, "repository.List")
	}
	for query.MapScan(m) {
		clientList = append(clientList, &core.Client{
			IdMeeting: m["idmeeting"].(string),
			IdSession: m["idsession"].(int),
			IdUser: m["iduser"].(int),
		})
		m = map[string]interface{}{}
	}
	log.Println("DEBUG", " Repository.List.Len: ", len(clientList))
	return clientList, nil
}

func (repository *CassandraRepository) Store(client *core.Client) error {
	var insertQuery = "insert into transmission.clients (idMeeting, idUser, tokenData) values (?,?,?)";
	err := repository.session.Query(insertQuery, &client.IdMeeting, &client.IdUser, &client.Token).Exec()
	if err != nil {
		log.Println("ERROR", "Error guardando cliente", err)
		return errors.Wrap(err, "repository.Store")
	}
	return nil
}

func (repository *CassandraRepository) Delete(token string) error {
	var meetingKey string
	var userKey int
	var meetingQuery = "select idMeeting, idUser from transmission.clients where tokenData=? ALLOW FILTERING"
	meetingOutput := repository.session.Query(meetingQuery, token).Iter()
	meetingOutput.Scan(&meetingKey, &userKey)

	var deleteQuery = "delete from transmission.clients where idMeeting=? and idUser=?";
	err := repository.session.Query(deleteQuery, &meetingKey, &userKey).Exec()
	if err != nil {
		log.Println("ERROR", "Error eliminando", err)
		return errors.Wrap(err, "repository.Delete")
	}
	return nil
}

func (repository *CassandraRepository) Update(token string, idSession int, media []byte) error {
	var updateQuery string
	var err error

	var meetingKey string
	var userKey int
	var meetingQuery = "select idMeeting, idUser from transmission.clients where tokenData=? ALLOW FILTERING"
	meetingOutput := repository.session.Query(meetingQuery, token).Iter()
	meetingOutput.Scan(&meetingKey, &userKey)

	if idSession < 0 {
		updateQuery = "update transmission.clients set idSession = ? where idMeeting=? and idUser=?";
		err = repository.session.Query(updateQuery, idSession, &meetingKey, &userKey).Exec()
	} else {
		updateQuery = "update transmission.clients set media = ? where idMeeting=? and idUser=?";
		err = repository.session.Query(updateQuery, media, &meetingKey, &userKey).Exec()
	}
	if err != nil {
		log.Println("ERROR", "Error actualizando", err)
		return errors.Wrap(err, "repository.Delete")
	}
	return nil
}

func (repository *CassandraRepository) Validate(token string) bool{
	var existToken bool
	var existsQuery = "select true from transmission.clients where tokenData=? ALLOW FILTERING"
	meetingOutput := repository.session.Query(existsQuery, token).Iter()
	meetingOutput.Scan(&existToken)
	return existToken == true
}
