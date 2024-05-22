# Container image that runs your code
FROM nginx

# Copies your code file from your action repository to the nginx folder
COPY index.html /usr/share/nginx/html

# Code file to execute when the docker container starts up (`entrypoint.sh`)
# ENTRYPOINT ["/entrypoint.sh"]
