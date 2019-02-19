FROM labstack/armor:0.4.11

COPY armor.yaml etc/armor/config.yaml
COPY public/ /var/www/

CMD ["-c", "/etc/armor/config.yaml"]