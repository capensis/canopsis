FROM node:8.12-alpine

WORKDIR /canopsis-next
ADD ./ /canopsis-next/

RUN cd /canopsis-next/ && \
    yarn install

EXPOSE 8080
CMD ["yarn", "run", "serve"]
