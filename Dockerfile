FROM scratch
ADD /yxkh //
ADD /config.json //
EXPOSE 8080
ENTRYPOINT [ "/yxkh" ]