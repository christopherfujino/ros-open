# bookworm, last updated 2026-02-28
FROM debian:13

HEALTHCHECK --interval=5m --timeout=2s \
  CMD curl --fail http://localhost/ || exit 1

# mailcap is for /etc/mime.types
RUN apt-get update \
  && apt-get upgrade -y \
  && apt-get install --no-install-recommends -y \
    cgit \
    lighttpd \
    mailcap \
  && apt-get clean

# Convenience dev deps
#RUN apt-get install --no-install-recommends -y \
#  vim \
#  man-db \
#  less \
#  curl

COPY ./lighttpd.conf /etc/lighttpd/conf.d/cgit.conf
COPY ./cgitrc /etc/cgitrc

RUN echo 'include "conf.d/cgit.conf"' >> /etc/lighttpd/lighttpd.conf

COPY ./entrypoint.sh /entrypoint.sh
COPY ./index.html /var/www/html/index.html
RUN rm /var/www/html/index.lighttpd.html

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
