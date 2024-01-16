# Lab service
This repository was created to implement a real workflow for the wards.

# Task
Create server to store and get inodes.
Inodes is linux file system entity that store metadata about file: [link](https://www.bluematador.com/blog/what-is-an-inode-and-what-are-they-used-for)

### Endpoints
- /api/v1/user [POST] to save inode
- /api/v1/user [GET] to get inode

### Index Node (inode) entity should contain the following fields

```JSON
{
  "id": 4242239,
  "file_name": "string",
  "filepath": "/etc/msql/conf.d",
  "file_type": "d",
  "mode" : "rwxrw-r--",
  "uid": 1,
  "gid": 1,
  "num_of_blocks": 4096,
  "size": 0,
  "timestamp": ""
}
```
Note: size = num_of_blocks * block_size

### Need to store in runtime
* using slice
* using hash map