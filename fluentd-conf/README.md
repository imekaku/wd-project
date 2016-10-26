### How to start fluentd
/opt/td-agent/embedded/bin/fluentd -c client.conf
/opt/td-agent/embedded/bin/fluentd -c server.conf

### Put fluentd log in file && Set log level warn
fluentd -c client.conf -q -o /home/work/fluentd-log/own-log/fluentd

[logging fluentd](http://docs.fluentd.org/articles/logging)

### Need plugins
```shell
# plugin that client need
/opt/td-agent/embedded/bin/fluent-gem install fluent-plugin-rewrite-tag-filter
/opt/td-agent/embedded/bin/fluent-gem install fluent-plugin-grep
/opt/td-agent/embedded/bin/fluent-gem install fluent-plugin-record-reformer
/opt/td-agent/embedded/bin/fluent-gem install fluent-plugin-forest
 
# plugin that client need
/opt/td-agent/embedded/bin/fluent-gem install fluent-plugin-forest
```
