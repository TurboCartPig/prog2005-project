[![pipeline status](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021-workspace/lindtvedtsebastian/cloud-project/developer-bot/badges/master/pipeline.svg)](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021-workspace/lindtvedtsebastian/cloud-project/developer-bot/-/commits/master)

# Developer bot

This is a Discord bot with two main pieces of functionality. It can extract deadlines from GitLab issue tracker. And it can create decision polls on special issues in GitLab. It does this by using GitLab's webhook interface, in order to receive data from GitLab, and simultaneously runs a discord bot, that facilitates user interaction.

Deadlines:
> TODO: Specify how deadlines work and how to interact with them

Decisions:
> TODO: Specify how decisions work and how to interact with them

# Development

The development of this project is centered around docker-compose and containerization. You *can* build the server as a binary manually, but you will have a more difficult time deploying it.

Both the docker-compose setup and a manual setup requires the presence of a `service-account-key.json` from GCP. This is used both to authorize Firestore and Google Cloud's Secret Manager, the latter is where the Discord bot token is stored and accessed securely. The account key requires the roles `INSERT FIRESTORE ADMIN ROLE` and `INSERT SECRET MANAGER SECRET ACCESSOR ROLE`. The key is discovered from the environment variable `GOOGLE_APPLICATION_CREDENTIALS`, which contains a path pointing to the service key, default should be `./service-account-key.json`. This is made available to the container through docker-compose's secrets or GCP's provisioning.

## Building

This project can be built in two ways.

### Manually

To build:
```bash
go build -o bin/developer-bot
```

To run:
```bash
./bin/developer-bot
```

Or build and run in one with:
```bash
go run
```

### With Docker

To build an image:
```bash
docker-compose build
```

To run a container based on built image
```bash
docker-compose up -d
```

You can also build and run the container all in one go:
```bash
docker-compose up -d --build
```

And to shut the container down:
```bash
docker-compose down
```

## CI

This project uses a GitLab's CI pipelines to validate the state of the project and to produce artifacts. The CI consists of multiple stages, the descriptions of which you can find below. The purpose of the CI is to catch errors as early as possible, to do regression testing, and to prevent faulty code from going into production. The CI as configured caches all that it can quite aggressively, in order to increase performance and save on network bandwidth. This is necessary as the runner is hosted on Skyhigh, and is quite resource constrained.

1. Build stage - Builds the project as a normal Go binary, just to make sure there are no config issues or trivial mistakes in syntax.
2. Lint stage - More detailed analysis of the code to find common mistakes and anti-patterns.

Stages yet to be added:
3. Docker build stage - Build the project as a Docker container, via Dockerfile or docker-compose.
4. Deployment stage - Deploy the docker image built in the previous stage.

# Deployment

There are a few considerations to take into account when considering how to deploy the bot. Mainly how it's built, how to run it and keep it running, and how to configure access to the secret token.

## GCP via Container optimized OS

This is the config the current deployment uses.

- Create a new service account with the following roles: ...
- Build the docker image locally and upload it to your preferred container registry
- Create a new VM instance in GCP
- Select deploy container
- Fill in the URL to your container image
- Select a container optimized OS
- Select the service account you just created
- Select allow http traffic
- Click start VM. Everything is automatic from here

## Manually on OpenStack or other IaaS solution

- Manually copy over a service account key
- Deploy using docker-compose

## Heroku

The current stack can not be deployed to heroku, due to heroku's poor handling of secrets. These were the instructions back when it worked.

There are two main methods for deploying the app to Heroku. The default Heroku way, using the heroku-20 stack. Or by deploying the app as a container using the container stack. The default way is configured by the `Procfile`, and the container way is configured by `heroko.yml` and the `Dockerfile`.

Deploy the default way:
```bash
# Login to heroku if you aren't already
heroku login
# Create the heroku app
heroku create
# Set the secret bot token
heroku config:set TOKEN=<...>
# Deploy to heroku
git push heroku main
```

To configure heroku to deploy the app as a container, run the following commands just before first deploying to heroku:
```bash
# Set the stack to container as opposed to heroku-20
heroku stack:set container
# Redeploy
git push heroku main
```

## Security

Quite a lot of time was spent figuring out how to deploy, and preferably develop, this project with security in mind. We considered any security through obscurity to be unacceptable. So hard coding or environment variables are out the door. We explored using HashiCorp Vault and many other similar products, before finally settling on using GCP's built in security features. Namely service accounts and the secrets manager. Here we set up granular permissions via service accounts. And securely store the discord bot token with encryption at rest and in transit. The service account key necessary to facilitate authentication and secure communicating is passed into the container with use of docker's built in secrets mechanisms, and works the same way when developing, as when deploying. This setup provides both a good level of security, and is relatively comfortable while developing.
