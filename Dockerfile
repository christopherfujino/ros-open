# bookworm, last updated 2026-02-28
FROM debian:13

# mailcap is for /etc/mime.types
RUN apt-get update \
  && apt-get upgrade -y \
  && apt-get install --no-install-recommends -y \
    cgit \
    lighttpd \
    mailcap \
  && apt-get clean

# Convenience dev deps
RUN apt-get install --no-install-recommends -y \
  vim

COPY ./lighttpd.conf /etc/lighttpd/conf.d/cgit.conf

RUN echo 'include "conf.d/cgit.conf"' >> /etc/lighttpd/lighttpd.conf

COPY ./entrypoint.sh /entrypoint.sh

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
#CMD ["lighttpd", "-D", "-f", "/etc/lighttpd/lighttpd.conf"]
