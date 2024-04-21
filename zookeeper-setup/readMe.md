**Pre-requisites**

1) Create a zoo.cfg file inc ase you want to set up your own configuration.
2) Change the data directory path in the conf file as per your data directory

---------------------------------------------------------------------------------------

**Deployment Steps:**

Step 1: Create a docker compose file with settings and volumes defined in it.
        Refer ./zoo.cfg file respective directory

Step 2: Run the docker cotainer
        docker compose up  (without detatched mode)
        docker compose up -d (with detatached mode)

Step 3: Enter the docker container
        docker exec -it <container name> bash

Step 4: Start Zookeper CLI
        bin/zkCli.sh

From here you can run the respective zookeeper commands 
