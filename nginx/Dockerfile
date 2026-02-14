FROM ubuntu:22.04
RUN apt-get update && apt-get install -y nginx openssl && rm -rf /var/lib/apt/lists/*
COPY nginx.conf /etc/nginx/nginx.conf
COPY default-*.conf /etc/nginx/
COPY index.html /usr/share/nginx/html
WORKDIR /app
COPY generate_cert.sh /app/generate_cert.sh
RUN chmod +x /app/generate_cert.sh && /app/generate_cert.sh
EXPOSE 80 8080 443
CMD ["nginx", "-g", "daemon off;"]