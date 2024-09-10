1. just like chord. Use a int list/slice for simulate the ring. When recieve the key -> get a hash in a range 2^32, and then find the closest node.
2. for dealing with situation that real node cannot evenly distribute in a ring, use virtual node to make them evenly distribute.
3. maybe add delete node later
4. also notice that add new node to the ring is aiming for make less node group change. ex: one node handle 10 keys, and new second node join, we want only five of them change group. 