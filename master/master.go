package master

import (
    . "SDSs/utils"
)

func StartMaster() {
    listener,err := net.Listen("tcp",MasterAddr)
    if err!= nil{
        fmt.Println("err = ",err)
        return
    }
    defer listener.Close()
    clientmanager = &ClientManager{}
    workermanager = &WorkerManager{}
    clientmanager.run()
    workermanager.run()
}
