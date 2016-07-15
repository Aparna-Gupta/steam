# Installing, Starting, and Using Steam and the H2O Scoring Service

This document is meant for H2O developers and describes how to install, start, and use Steam on an in-house YARN cluster. External users should review the [README.md](docs/README.md) file within the **/docs** folder.


## Requirements
- Web browser and an Internet connection
- steamY repository
- Go (available from <a href="https://golang.org">golang.org</a>) 
- Access to the steamY repository
- SSH access to a Jetty server running YARN
- Typescript
- Node.js
- JDK 1.7 or greater
- H2O AutoML for Apache HDP2.2 or CDH 5.5.3 (internal only)
- postgresql for using the CLI

## Building the H2O Scoring Service

The H2O Scoring Service Builder is an application that allows you to perform the following through either a web UI or command line:

- Compile a POJO, and then build a Jar file from a POJO and a gen-model file
- Compile the POJO, and then build a War file that is a service from a POJO, gen-model. You can then run the war file in Jetty, Tomcat, etc.
- Build a War file with a POJO predictor and Python pre-preprocessing

Perform the following steps to build the H2O Scoring Service Builder:

1. In a terminal window, navigate to the **steamY/scoring-service-builder** folder.

2. Run `./gradlew build` to build the service.

3. You will see a **BUILD SUCCESSFUL** message upon completion. Run 	`./gradlew jettyRunWar` to run the builder service.

4. Open a browser and navigate to localhost:55000 to begin using the H2O Scoring Service Builder web UI. 

### Testing the Scoring Service Builder

**Using the Web UI**

When the Builder Service is running, you can make a War file using the Web UI.

The following screenshot shows how to make a War file using a POJO file and a Jar file. Note that these files are included in the  **steamY/scoring-service-builder/examples/example-pojo** folder. 

![Make War](scoring-service-builder/images/make_war.png)


**Using the CLI**

Note that when the Builder Service is running, you can also make a war file using command line arguments. For example:

		curl -X POST --form pojo=@examples/example-pojo/gbm_3f258f27_f0ad_4520_b6a5_3d2bb4a9b0ff.java --form jar=@examples/example-pojo/h2o-genmodel.jar localhost:55000/makewar > example.war

 where:
 
 - `gbm_3f258f27_f0ad_4520_b6a5_3d2bb4a9b0ff.java` is the POJO file from H2O
 - `h2o-genmodel.jar` is the corresponding Jar file from your version of H2O

Both of the above files are included in the **steamY/scoring-service-builder/examples/example-pojo** folder.

#### Importing a War File

When the H2O Scoring Service Builder is up and running, open another terminal window, navigate to the **steamY/scoring-service-builder** folder, and run the following command to import the war file into the H2O Prediction service. The example below uses the **example.war** file that was built using the CLI in the previous section.

		java -jar jetty-runner-9.3.9.M1.jar --port 55001 example.war

This starts the H2O Prediction Service at localhost:55001. You can view this web service at http://localhost:55001.

![Example Service](scoring-service-builder/images/example_service.png)

Close the browser and the terminal window when the testing is complete.


## Steam Installation
Perform the following steps to install Steam.

1. Open a Terminal window and install Go (available from <a href="https://golang.org" target="_blank">golang.org</a>).

2. Export the GOPATH in the location where Steam will reside. This can be added to your bash profile or specified each time you run Steam. For example:

	`export GOPATH=$HOME/Desktop` 

3. Create the directory where Steam will reside. This must be in GOPATH location. For example:

	`mkdir -p $GOPATH/src/github.com/h2oai`
	
	>**Note**: The path MUST include `/src/github.com/h2oai`. 

4. Change directories to the new **h2oai** folder, and clone the repository in this folder. Enter your git username and password when prompted. 

		cd Desktop/src/github.com/h2oai
		git clone https://github.com/h2oai/steamY
	
5. Change directories to **steamY** and then run `make linux`.

		cd steamY
		make
		make linux

	>***Note***: If you are building this on your own local machine, then you can run `make` instead of `make linux`. This will only work if your local machine has YARN. 

You will see a `BUILD SUCCESSFUL` message when the installation is complete. At this point, you are ready to start using Steam. 


## Starting Steam
When Steam is installed, perform the following steps to start and use Steam. Note that two terminal windows will remain open: one for the Jetty server and one for Steam.

1. Open a terminal window. In the **steamY** directory, copy the **automl-hdp2.2.jar** file to your local machine. This engine is required in order to run any AutoML jobs in Steam. 

		scp <user>@<domain>:./automl-hdp2.2.jar .

2. ssh to the machine running YARN, changing `<user>` and `<domain>` below with appropriate values. Specify your password when prompted. 

	`ssh <user>@<domain>`
	
	The machine will include the following folder and files:
	
	- steam--linux-amd64.tar.gz
	- automl-hdp2.2.jar	
	
3. On the YARN machine, change directories to the Steam directory, then run the H2O Java Model Compilation Service. This service is required for Steam to be able to compile and build models for deployment purposes:

		cd steam--linux-amd64/steam
		java -jar var/master/assets/jetty-runner.jar --port 8811 var/master/assets/ROOT.war

5. Open another terminal window and again, ssh to the machine running YARN, changing `<user>` and `<domain>` below with appropriate values. Specify your password when prompted. 

	`ssh <user>@<domain>` 
	
6. Untar the **steam--linux-amd64.tar.gz** package.

	`tar -xzf steam--linux-amd64.tar.gz`

6. Change directories to the new **steam--linux-amd64** folder, then start the Steam master node. For example, the following commands will start Steam on localhost.

		cd steam--linux-amd64
		./steam serve master
		
	>**Note**: You can view all available options for starting Steam on the master using `./steam help serve master`

	You will see a message similar to the following when Steam starts successfully.

		2016/04/28 13:34:56 steam v build 2016-04-28T20:15:00+0000
		2016/04/28 13:34:56 Working directory: /home/seb/steam--linux-amd64/var/master
		2016/04/28 13:34:56 WWW root: /home/seb/steam--linux-amd64/var/master/www
		2016/04/28 13:34:57 Priming datastore for first time use...
		2016/04/28 13:34:57 Datastore location: /home/seb/steam--linux-amd64/var/master/db/steam.db
		2016/04/28 13:34:57 Web server listening at :9000
		2016/04/28 13:34:57 Point your web browser to http://localhost:9000/

## Using Steam

In a Web browser, navigate to the Steam Web server (for example, http://172.16.2.182:9000).

### Adding an Engine
An empty Steam UI will display. Before performing any tasks, you must first add an Asset. 

1. Click the **Assets** icon (<img src="docs/images/icon_assets.png" alt="Thumbnail" style="width: 25px;" />) on the left navigation panel, then select **Add Engine**. 

	![](docs/images/add_engine.png)

2. Browse to the **automl-hdp2.2.jar** file on your local machine, then click **Upload**. 

	>***Note***: This file was added during the "Install Steam" steps. 

### Starting a Cloud

Clouds can be configured after the engine asset was successfully added. 

1.  Click the **Clouds** icon (<img src="docs/images/icon_clouds.png" alt="Thumbnail" style="width: 25px;" />) on the left navigation panel, then select **Start a Cloud**. 

2. Enter/specify the following information to set up your cloud:

	a. A name for the cloud

	b. The version of H2O that will run on the cloud.

	c. The number of nodes on the cloud.
	
	d. The amount of memory available on each node. Be sure to include the unit ("m" or "g").
	
	![](docs/images/add_cloud.png)
	
3. Click **Start Cloud** when you are finished.

>***Note***: You can view a stream of the cloud creation log in the terminal window that is running Steam. In the UI, Steam will respond with an error if the cloud configuration is incorrect (for example, if you specify more nodes than available on the cluster). 

The Cloud Details page opens upon successful completion. This page shows the cloud configuration information and includes a link to the H2O Flow URL. From this page, you can begin building model. 

![](docs/images/cloud_details.png)

### Adding a Model
Models are created from the Cloud Details page. When building a model, you will need to provide the location of the dataset that you will use as well as the response column. 

1. Click the **Build a Model** button on the bottom of the Cloud Details page.

2. Enter a path for the dataset that you want to use to build the model. 

>***Note***: If you choose to use a local dataset, then that dataset must reside in the same folder/path on each node in the cluster.

3. Specify the column that will be used as the response column in the model. 

4. Specify the maximum run time in seconds. H2O will return an error if the model build stalls after this threshold is reached.

5. Click **Start Building** when you are finished. 

	![](docs/images/build_model.png)

### Viewing Models	

Click the **Models** icon (<img src="docs/images/icon_models.png" alt="Thumbnail" style="width: 25px;" />) on the left navigation panel to view models that were successfully created. 

These models are processed using H2O's AutoML algorithm, which determines the best method to use to build the model. The model name includes this method. So, for example, if Steam returns a model named "DRF_model...", then this indicates that DRF was the algorithm that provided the best result.

### Deploying Models

After a model is built, the next step is to deploy the model in order to make/view predictions.

1. Click the **Models** icon (<img src="docs/images/icon_models.png" alt="Thumbnail" style="width: 25px;" />) on the left navigation panel.

2. Select the model that you want to use, then click the **Deploy this Model** button on the bottom of the page.

3. Specify the port to use for the scoring service. 

	>***Note***: Steam will return an error if you specify a port that is already being used.

4. Click **Deploy** when you are done.

### Making Predictions

Successfully deployed models can be viewed on the **Services** page. On this page, click the **Endpoint** link to open the H2O Prediction Service and begin making predictions.  

## Using Steam with Flow

As with other H2O products, Flow can be used alongside Steam when performing machine learning tasks.

On the Cloud Details page, click the Address link to open H2O Flow in a new tab. 

![](docs/images/h2o_flow.png)

>***Note***: Refer to the H2O Flow documentation for information on how to use Flow. 

## Stopping Steam	

When you are finished, use the following process to safely shut down Steam:

1. On the Services page, stop all running services.

2. Stop all running clouds.


## User Management

This section describes the user management features available in H2O Steam when Steam is connected to a database. The Steam database is useful for creating a new bash script and for automatic model deployment. The Steam database supports setup via Python functions and CLI commands. 

This section includes the following subsections:

- [Terms](#terms)
- [Privileges/Access Control](#privileges)
- [Authorization](#authorization)
- [User Management Setup](#user management setup)
- [User Management Workflow](#user management workflow)
- [CLI Command Reference](#CLI Command Reference)

<a name="terms"></a>
### Terms

The following lists common terms used when describing Steam User Management.  

- **Entities** represent *objects* in Steam. Examples of entities include Roles, Workgroups, Identities, Clusters, Projects, Models, and Services (engines). 

- **Identities** represent *users* in Steam. Users sign in using an Identity, and then perform operations in Steam.

- **Permissions** determine what operations you can perform. Examples of permissions include *Manage Clusters*, *View Clusters*, *Manage Models*, *View Models*, and so on.

- **Privileges** determine the entities that you can perform operations on (i.e., data / access control).


<a name="privileges"></a>
### Privileges/Access Control

Privileges are uniquely identified by the entity in question and the kind of privilege you have on the entity.

The following privileges are available on an entity:

- **Own** privileges allow you to share, view, edit, and delete entities.

- **Edit** privileges allow you to view and edit entities, but not share or delete them.

- **View** privileges allow you to view entities, but not share, edit, or delete them.

When you create an entity, you immediately *Own* it. You can then share this entity with others and award them either *Edit* or *View* privileges. Entities are allowed to have more than one owner, so you can also add additional owners to entities. 

The following table lists the kind of privileges you need in order to perform specific operations on entities:


        Entity               Own  Edit View
        -----------------------------------
        Role
          Read               x    x    x
          Update             x    x
          Assign Permission  x    x
          Delete             x
          Share              x
          
        Workgroup
          Read               x    x    x
          Update             x    x
          Delete             x
          Share              x
        
        Identity
          Read               x    x    x
          Assign Role        x    x
          Assign Workgroup   x    x
          Update             x    x
          Delete             x
          Share              x
        
        Cluster
          Read               x    x    x
          Start/Stop         x
        
        Project
          Read               x    x    x
          Assign Model       x    x
          Update             x    x
          Delete             x
          Share              x
        
        Engine, Model
          Read               x    x    x
          Update             x    x
          Delete             x
          Share              x

        

<a name="authorization"></a>
### Authorization

Permissions and privileges are set up using Roles and Workgroups, respectively.

- Identities cannot be linked directly to permissions. For that, you'll need Roles.

- Identities cannot be linked directly to privileges on entities. For that, you'll need Workgroups, i.e. when you share entities with others, you would be sharing those entities with workgroups, not individuals.

#### Roles
A **Role** is a named set of permissions. Roles allow you define a cohesive set of permissions into operational roles and then have multiple identities *play* those roles, regardless of access control.
For example:

- a *Data Scientist* role can be composed of the permissions *View Clusters*, *Manage Models*, *View Models*.
- an *Operations* role can be composed of the permissions *View Models*, *View Services*, *Manage Services*,
- a *Manager* role can be composed of the permissions *Manage Roles*, *View Roles*, *Manage Workgroups*, *View Workgroups*

#### Workgroups
A **Workgroup** is a named set of identities. Workgroups allow you to form collections of identities for access control purposes. For example, a *Demand Forecasting* workgroup can be composed of all the users working on demand forecasting, regardless of their role. This workgroup can be then used to control access to all the clusters, projects, models and services that are used for demand forecasting. 


<a name="user management setup"></a>
### User Management Setup

This section describes how to set up and start the Steam  CLI for user management. Five terminal windows will be open the first time you run this setup; four terminal windows will be open for subsequent logins.

1. Open a terminal window and start postgresql. This should be started from the folder where posgresql was installed.

		postgres -D /usr/local/var/postgres

2. Open a second terminal window to create a new user for the Steam database and then create the database. The commands below only need to be performed once. The example below creates a steam **superuser** with a password ``st3amUser`` before creating the Steam database. Be sure to provide a secure password, and be sure to remember the password that you enter. This will be required each time you log in to Steam. 

		createuser -P steam 
		Enter password for new role: st3amUser
		Enter it again: st3amUser
		# Change directories to the Steam /var/master/scripts folder.
		cd steam-master-darwin-amd64/var/master/scripts
		./create-database.sh

3. Open a third terminal window. Navigate to the Steam folder and run the following command from within the Steam folder to start the Steam compilation service. 

		cd ../../../
		java -jar var/master/assets/jetty-runner.jar var/master/assets/ROOT.war

4. Open a fourth terminal window. From within the Steam folder, start Steam using the password that you provided in Step 2. This starts Steam on localhost:9000.

		./steam serve master --superuser-name=superuser --superuser-password=superuser
		
5. <a name="step5"></a>Open a fifth terminal window. From within the Steam folder, log in to the maching running Steam (localhost:9000). Use the password that you provided in Step 2.

		./steam login localhost:9000 --username=superuser --password=superuser

6. Run the following to verify that the CLI is working correctly.

		./steam help
		
	You should see the following output.
	
		steam vmaster build 2016-07-01T16:26:15+0000: Command Line Interface to Steam
		
		Usage:
		  steam [command]
		
		Available Commands:
		  login       Sign in to a Steam server.
		  reset       Reset Steam client configuration.
		  serve       Launch a new service.
		  start       Start a new resource.
		  stop        Stop the specified resource.
		  register    Register an external resource.
		  unregister  Unregister an external resource.
		  deploy      Deploy a resource of the specified type.
		  get         List or view resources of the specified type.
		  delete      Deletes the specified resource from the database.
		  import      Import a resource of the specified type into steam.
		  create      Creates an instance of the specified resource.
		  deactivate  Deactivate and entity type.
		  update      Updates an entity in the database.
		  link        Add authentication permissions.
		  unlink      Remove authentication permissions.

		Flags:
		  -v, --verbose[=false]: verbose output

		Use "steam [command] --help" for more information about a command. 


<a name="user management workflow"></a>
### User Management Workflow

The steps below provide a common workflow to follow when creating users. This workflow is followed in the example that follows.

1. Define roles based on operational needs.
2. Define workgroups based on data / access control needs.
3. Then add a new user:

 -	Create the user's identity.
 - Associate the user with one or more roles.
 - Optionally, associate the user with one or more workgroups. 

#### Example

The following example creates sample roles, workgroups, and users using the CLI. Refer to the [CLI Command Reference](#CLI Command Reference) section for information about all of the commands available in the CLI. These commands are run from the terminal window used to log in to Steam ([Step 5](#step5) above).

		# Create engineer role and link that role to permissions
		./steam create role engineer --desc="a default engineer role"
		./steam link role engineer ViewClusters ViewModels ViewWorkgroups
		
		# Create data scientist role and link that role to permissions
		./steam create role datascience --desc="a default data scientist role"
		./steam link role datascience ViewClusters CreateModels ViewWorkgroups
		
		# Create preparation and production workgroups
		./steam create workgroup preparation --desc="data prep group"
		./steam create workgroup production --desc="production group"
		
		# Create two users - Bob and Jim
		./steam create identity bob bobSpassword
		./steam create identity jim j1mSpassword
		
		# Link Bob to engineer role; link Jim to datascience role
		./steam link identity bob role engineer
		./steam link identity jim role datascience
		
		# Link Bob to preparation workgroup; link Jim to production workgroup
		./steam link identity bob workgroup preparation
		./steam link identity jim workgroup production


### Stopping the Steam Database

Use Ctrl+C in each of the Steam, Compilation Service, and postgres terminal windows to stop the services end your session. 

<a name="CLI Command Reference"></a>
# CLI Command Reference

- [`create identity`](#create identity)
- [`deactivate identity`](#deactivate identity)
- [`delete cluster`](#delete cluster)
- [`delete engine`](#delete engine)
- [`delete model`](#delete model)
- [`delete role`](#delete role)
- [`delete service`](#delete service)
- [`delete workgroup`](#delete workgroup)
- [`deploy engine`](#deploy engine)
- [`get clusters`](#get clusters)
- [`get engines`](#get engines)
- [`register cluster`](#register cluster)
- [`start cluster`](#start cluster)
- [`stop cluster`](#stop cluster)
- [`unregister cluster`](#unregister cluster)


------

#### <a name="create identity"></a>`create identity`

**Description**

Creates an instance of the specified resource.

**Usage**

	./steam create identity [username] [password]


**Parameters**

- `[username]`: Enter a string for the new user name
- `[password]`: Enter a string for the new user's password

**Example**

The following example creates a new user with a username of "minky" and a password of "m1n5kypassword". 
 
	./steam create minsky m1n5kypassword
	
------

#### <a name="deactivate identity"></a>`deactivate identity`

**Description**

Deactivates an identity based on the specified username.

**Usage**

	./steam deactivate identity [username]

**Parameters**

- `[username]`: Specify the username of the identity that you want to deactivate.

**Example**

The following example deactivates user "minsky". 

	./steam deactivate minsky 

-----

#### <a name="delete cluster"></a>`delete cluster`

**Description**

Deletes the specified cluster from the database.

**Usage**


**Parameters**


**Example**

-----

#### <a name="delete engine"></a>`delete engine`

**Description**

Deletes the specified engine from the database.

**Usage**


**Parameters**


**Example**

-----
#### <a name="delete model"></a>`delete model`

**Description**

Deletes the specified model from the database.

**Usage**


**Parameters**


**Example**

-----

#### <a name="delete role"></a>`delete role`

**Description**

Deletes the specified role from the database.

**Usage**


**Parameters**


**Example**

-----

#### <a name="delete service"></a>`delete service`

**Description**

Deletes the specified service from the database.

**Usage**


**Parameters**


**Example**

-----

#### <a name="delete workgroup"></a>`delete workgroup`

**Description**

Deletes the specified workgroup from the database.

**Usage**


**Parameters**


**Example**

-----

#### <a name="deploy engine"></a>`deploy engine` 

**Description**

Deploys an H2O engine. After an engine is successfully deployed, it can be specified when starting a cluster. (See [`start cluster`](#start cluster).) 

**Usage**

	./steam deploy engine [path/to/engine]

**Parameters**

- `[path/to/engine]`: Specify the location of the engine that you want to deploy. 

**Example**

The following specifies to deploy the H2O AutoML engine.

	./steam deploy engine ../engines/automl-hdp2.2.jar

-----

#### <a name="get clusters"></a>`get clusters`

**Description** 

Retrieves a list of clusters.

**Usage**

	./steam get clusters

**Parameters**

None

**Example**

The following example retrieves a list of clusters that are running H2O and are registered in Steam. (See [`register cluster`](#register cluster).)

	./steam get clusters
	NAME		ID	ADDRESS			STATE	TYPE		AGE
	user     	1	localhost:54321	started	external	2016-07-01 11:45:58 -0700 PDT

-----

#### <a name="get engines"></a>`get engines`

**Description** 

Retrieves a list of deployed engines.

**Usage**

	./steam get engines

**Parameters**

None

**Example**

The following example retrieves a list of engines that have been deployed. (Refer to [`deploy engine`](#deploy engine).)

	./steam get engines
	NAME			ID	AGE
	h2o-genmodel.jar	1	2016-07-01 13:30:50 -0700 PDT
	h2o.jar			2	2016-07-01 13:32:10 -0700 PDT

-----


#### <a name="register cluster"></a>`register cluster`

**Description**

Registers a cluster that is currently running H2O (typically a local cluster). Once registered, the cluster can be used to perform machine learning tasks through Python, R, and Flow. The cluster will also be visible in the Steam web UI. 

Note that clusters that are started using this command can be stopped from within the web UI or using [`unregister cluster`](#unregister cluster). You will receive an error if you attemt to stop registered clusters using the `stop cluster` command. 

**Usage**

	./steam register cluster [address]

**Parameters**

- `[address]`: Specify the IP address and port of the cluster that you want to register.

**Example**

The following example registers Steam on localhost:54321. Note that this will only be successful if H2O is already running on this cluster. 

	./steam register cluster localhost:54321
	Successfully connected to cluster 2 at address localhost:54321

-----

#### <a name="start cluster"></a>`start cluster`

**Description**

Starts a new cluster through YARN using a specified engine. Note that this command is only valid when starting Steam on a YARN cluster. To start Steam on a local cluster, use [`register cluster`](#register cluster) instead.

**Usage**

	./steam start cluster [id] [engineid] --size=[numNodes] --memory=[string]

**Parameters**

- `[id]`: Enter an ID for this new cluster.
- `[engineid]`: Specify the ID of the engine that this cluster will use. If necessary, use [`get engines`](#get engines) to retrieve a list of all available engines.
- `--size=[numNodes]`: Specify an integer for the number of nodes in this cluster.
- `--memory=[string]`: Enter a string specifying the amount of memory available to Steam in each node (for example, "1024m", "2g", etc.)

**Example**

The following example retrieves a list of engines, then starts a cluster through YARN using one from the list. The cluster is configured with 2 nodes that are 2 gigabytes each. 

	./steam get engines
	NAME				ID	AGE
	h2o-genmodel.jar	1	2016-07-01 13:30:50 -0700 PDT
	h2o.jar			2	2016-07-01 13:32:10 -0700 PDT
	./steam start cluster 9 1 --size=2 --memory=2g
	
-----

#### <a name="stop cluster"></a>`stop cluster`

**Description**

Stops a YARN cluster that was started through the CLI or web UI. (See [`start cluster`](#start cluster).) Note that you will receive an error if you attempt to stop a cluster that was started using `register cluster`. 

**Usage**

	./steam stop cluster [id] 

**Parameters**

- `[id]`: Specify the ID of the cluster that you want to stop. If necessary, use [`get clusters`](#get clusters) to retrieve a list of clusters. 

**Example**

The following example stops a cluster that has an ID of 9.

	./steam stop cluster 9

-----

#### <a name="unregister cluster"></a>`unregister cluster`

**Description**

Stops a cluster that was registered through the CLI or the web UI. (See [`register cluster`](#register cluster).) Note that you will receive an error if you attempt to unregister a cluster that was started using `start cluster`. 

**Usage**

	./steam unregister cluster [id] 

**Parameters**

- `[id]`: Specify the ID of the cluster that you want to stop. If necessary, use [`get clusters`](#get clusters) to retrieve a list of clusters. 

**Example**

The following example stops a cluster that has an ID of 9. 

	./steam unregister cluster 2
	Successfully unregisted cluster %d 2

-----
