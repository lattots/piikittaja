FROM debian:trixie-slim

# Set environment variables to prevent interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Update the package list and install the MariaDB client package
# The --no-install-recommends flag keeps the image size small
RUN apt-get update \
	&& apt-get install -y --no-install-recommends mariadb-client \
	# Clean up the package cache to minimize the final image size
	&& rm -rf /var/lib/apt/lists/*

CMD ["mysql"]
