# PALD - port allocator daemon

`pald` is a very simple, intentionally feature-poor local daemon aimed to keep registry of port allocations on a local system. Due to the lack of time to address security concerns `pald` only listens on _localhost_ interfaces currently.

It is expected to operate on the port range of 49152-65535, as specified in [Section 6 of RFC-6335](http://tools.ietf.org/html/rfc6335#section-6). Any contiguous range of valid port numbers to be allocated can be specified in configuration.

## Configuration

Configuration is read from either system-wide and user-specific config._ext_ files. The state of assigned services persisted in the `dump` file in a location for user-specific configuration file. Configuratoin file location reflects executable name.

Dump file format is undecided yet and likely will be changed in the future.

Configuration file expected to be in the [TOML](https://github.com/toml-lang/toml) or othe formats as implemented by the `Viper` package used in `pald`. Here is what can be specified in the config file:

<table>
<tr><th>key</th><th>type</th><th>default</th><th>description</th></tr>
<tr><td>port_listen</td><td>uint16</td><td>49200</td><td>A port on which the <code>pald</code> process will listen for port queries and allocation requests</td></tr>
<tr><td>port_min</td><td>uint16</td><td>49201</td><td>The lowest (first) port available for allocation</td></tr>
<tr><td>port_max</td><td>uint16</td><td>49999</td><td>The highest (last) port available for allocation</td></tr>
<tr><td>dump_file</td><td>string</td><td>~/.pald/dump</td><td>The default dump file location where the service will persist the state while down</td></tr>
</table>

## HTTP interface

All requests are available as either HTTP GET or HTTP POST, e.g.

    http://localhost:49200/get?service=service-name

    REPLY=`curl -d service=service-name -o - -s -f http://localhost:49200/get`
    echo $?
    echo $REPLY

Three URLs are currently supported (with HTTP reply codes):

<table>

<tr><th>action</th><th>URL</th><th>param=value</th></tr>

<tr><td>Query</td><td>/get</td><td>service=name</td></tr><tr><td colspan="3" style="padding: 0.5em 0em 1.5em 2em;"><code>200</code> - a found port number<br />
<code>404</code> - an error message if there is no port registered with the requested service<br />
<code>400</code> - an error message in case of all other errors</td></tr>

<tr><td>Register</td><td>/set</td><td>service=name</td></tr><tr><td colspan="3" style="padding: 0.5em 0em 1.5em 2em;"><code>200</code> - an assigned port number<br />
<code>412</code> - registration failed because no more port numbers available in the configured range<br />
<code>400</code> - an error message in case of all other errors</td></tr>

<tr><td>Delete</td><td>/del</td><td>port=number</td></tr><tr><td colspan="3" style="padding: 0.5em 0em 1.5em 2em;"><code>200</code> - OK as a success indication (including port not found)<br />
<code>400</code> - an error message in case of all other errors</td></tr>
</table>

## Porting to other platforms

At this time `pald` is only compatible with Mac OS X, but it is easy to fix. Please, add an appropriate `internal\platform\specific_<platform>.go` file for your platform and send me a pull request.
