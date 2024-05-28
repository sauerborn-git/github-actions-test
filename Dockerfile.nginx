# Container image that runs your code
FROM nginx

# Copies your code file from your action repository to the nginx folder
COPY index.html /usr/share/nginx/html

# once everythign execute, exit cleanly to not start the nginx server
RUN exit 0

# Code file to execute when the docker container starts up (`entrypoint.sh`)
# ENTRYPOINT ["/entrypoint.sh"]
