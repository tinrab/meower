FROM scylladb/scylla:2.1.1

COPY create.cql entrypoint-wrap.sh /

ENTRYPOINT [ "/entrypoint-wrap.sh" ]
