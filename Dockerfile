FROM scratch
ADD /workdiary //
ADD /config.json //
EXPOSE 8080
ENTRYPOINT [ "/workdiary" ]