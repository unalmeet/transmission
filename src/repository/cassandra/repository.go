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
	CREATE_TABLE       = "create table if not exists transmission.clients (idMeeting text, idSession int, media blob, primary key (idMeeting, idSession));"
	KEYSPACE_NAME      = "transmission"
	CASSANDRA_URL      = "CASSANDRA_URL"
	CASSANDRA_USERNAME = "CASSANDRA_USERNAME"
	CASSANDRA_PASSWORD = "CASSANDRA_PASSWORD"
	PAGE_SIZE          = 10
)

type cassandraRepository struct {
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
	repo := &cassandraRepository{}
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

func (repository *cassandraRepository) List(idMeeting string) ([]*core.Client, error) {
	var clientList []*core.Client
	var selectQuery = "select * from transmission.clients where idMeeting=?"
	m := map[string]interface{}{}
	query := repository.session.Query(selectQuery, idMeeting).Iter()
	if query == nil {
		log.Println("ERROR", "Error obteniendo clientes")
		return nil, errors.Wrap(nil, "repository.List")
	}
	for query.MapScan(m) {
		clientList = append(clientList, &core.Client{
			IdMeeting: m["idmeeting"].(string),
			IdSession: m["idsession"].(int),
			Media:     m["media"].([]byte),
		})
		m = map[string]interface{}{}
	}
	log.Println("DEBUG", " Repository.List.Len: ", len(clientList))
	return clientList, nil
}

func (repository *cassandraRepository) Store(client *core.Client) error {
	var insertQuery = "insert into transmission.clients (idMeeting, idSession, media) values (?,?,?)";
	err := repository.session.Query(insertQuery, &client.IdMeeting, &client.IdSession, &client.Media).Exec()
	if err != nil {
		log.Println("ERROR", "Error guardando cliente", err)
		return errors.Wrap(err, "repository.Store")
	}
	return nil
}

func (repository *cassandraRepository) Delete(idMeeting, idSession string) error {
	var deleteQuery = "delete from transmission.clients where idMeeting=? and idSession=?";
	err := repository.session.Query(deleteQuery, idMeeting, idSession).Exec()
	if err != nil {
		log.Println("ERROR", "Error eliminando", err)
		return errors.Wrap(err, "repository.Delete")
	}
	return nil
}
