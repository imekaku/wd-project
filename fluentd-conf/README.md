## HOW TO START FLUENTD
/opt/td-agent/embedded/bin/fluentd -c client.conf
/opt/td-agent/embedded/bin/fluentd -c server.conf

## PUT FLUENTD LOG IN FILE & LOG LEVEL WARN
fluentd -c client.conf -q -o /home/work/fluentd-log/own-log/fluentd
[logging fluentd](http://docs.fluentd.org/articles/logging)

## NEED PLUGIN
```shell
# plugin that client need
/opt/td-agent/embedded/bin/fluent-gem install fluent-plugin-rewrite-tag-filter
/opt/td-agent/embedded/bin/fluent-gem install fluent-plugin-grep
/opt/td-agent/embedded/bin/fluent-gem install fluent-plugin-record-reformer
/opt/td-agent/embedded/bin/fluent-gem install fluent-plugin-forest
 
# plugin that client need
/opt/td-agent/embedded/bin/fluent-gem install fluent-plugin-forest
```
