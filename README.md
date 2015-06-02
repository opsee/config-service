# Config Service

The config service is a simple Go application that talks to Etcd in order to
manage versions of configuration information.

## SSH Tunnel

In order to use the config service, you need to setup an SSH tunnel to the
fleet you wish to deploy configurations. To do so, simply start an SSH 
connection:

`ssh -fN -L 2379:127.0.0.1:2379 cluster-address`

## Commands

### Get

The `get` command allows you to retrieve a version of a configuration
file for a named service.

```
λ cf get bartnet current
{
    "db-spec":
    {
        "classname":"org.postgresql.Driver",
        "subprotocol":"postgresql",
        "subname":"auth_test",
        "user":"cliff",
        "password":"",
        "max-conns":6,
        "min-conns":1,
        "init-conns":1
    },
    "secret":"abc123",
    "thread-util":0.9,
    "max-threads":64,
    "server":
    {
        "port":8080
    },
    "bastion-server":
    {
        "port":4080
    }
}
```

### Set

## Debugging

First, make sure you can talk to etcd:

```
λ go get github.com/coreos/etcd/etcdctl

λ go install github.com/coreos/etcd/etcdctl

λ etcdctl ls
/coreos.com
/opsee.co
```

If not, make sure your SSH tunnel is up and running.

If you do not see /opsee.co, then there is maybe a problem!

## TODO

Cross-AZ replication for Etcd so that writing configuration to one cluster 
sends that version of the configuration to every cluster.
