# Simple nginx image
FROM nginx

# WORKDIR /Users/cuilide/Desktop/blockchain app/web app/client
COPY . /usr/share/nginx/html
EXPOSE 80
RUN 'ls'

CMD ["nginx", "-g", "daemon off;"]

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "curl localhost:8080" ]