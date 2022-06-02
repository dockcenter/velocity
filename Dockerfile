FROM eclipse-temurin:17-jre-alpine

LABEL org.opencontainers.image.vendor="Dockcenter"
LABEL org.opencontainers.image.title="Velocity"
LABEL org.opencontainers.image.description="Automatically built Docker image for Velocity"
LABEL org.opencontainers.image.documentation="https://github.com/dockcenter/velocity/blob/main/README.md"
LABEL org.opencontainers.image.authors="Chao Tzu-Hsien <danny900714@gmail.com>"
LABEL org.opencontainers.image.licenses="MIT"

ENV JAVA_MEMORY="512M"
ENV JAVA_FLAGS="-XX:+UseStringDeduplication -XX:+UseG1GC -XX:G1HeapRegionSize=4M -XX:+UnlockExperimentalVMOptions -XX:+ParallelRefProcEnabled -XX:+AlwaysPreTouch"

WORKDIR /data

RUN addgroup -S velocity && \
    adduser -S velocity -G velocity && \
    chown velocity:velocity /data

USER velocity

VOLUME /data

EXPOSE 25577

COPY --chown=velocity velocity/velocity-*.jar /opt/velocity/velocity.jar

ENTRYPOINT java -Xms$JAVA_MEMORY -Xmx$JAVA_MEMORY $JAVA_FLAGS -jar /opt/velocity/velocity.jar