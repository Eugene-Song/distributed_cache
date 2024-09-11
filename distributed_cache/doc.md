#### Summary for building distributed cache
1. The general process
                               yes
recieve key --> if in local cache -----> return the cache data
                |  no                       yes  
                |-----> if get from peers -----> use http req to get from peer --> use consistent hash for pick nodes                                              yes
                            |                                                            |-----> is a remote? -------------------> HTTP req --> success-----> return response
                            |                                                                   |  no                yes                               no
                            |                                                                    |----------------------------------------------> go back to local
                            |  No
                            |-----> use call back, and add the data from data source to local cache --> return cacche

2. General three part
Get local cache:
The core is LRU cache. 

Get from node
When Multiple Nodes, need consistent hash for load balance and maintain nodes relationships. Peer.go is like the client when sending req for getting data from nodes. server.go will be like an interface for handling communication for recieving data and access the actual cache entry point.
(currently, also has function for setting up nodes, may clean up and extract for a new file. Two functions need to be cleaned. PickPeer and SetHttpPool)


Get from data source
Use byteview to simulate data source