+++
title = "Logger Plugin"
description = "Logger plugin logs HTTP requests"
[menu.main]
  name = "Logger"
  parent = "plugins"
  weight = 3
+++

Logs HTTP requests

## Configuration

<table>
  <thead>
    <tr>
      <th align="left">Name</th>
      <th align="left">Type</th>
      <th align="left">Value</th>
      <th align="left">Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="left"><code>name</code></td>
      <td align="left">string</td>
      <td align="left">logger</td>
      <td align="left">Plugin name</td>
    </tr>
    <tr>
      <td align="left"><code>format</code></td>
      <td align="left">string</td>
      <td align="left">
        <code>{"time":"${time_rfc3339_nano}",
        "remote_ip":"${remote_ip}",
        "host":"${host}",
        "method":"${method}",
        "uri":"${uri}",
        "status":${status},
        "latency":${latency},
        "latency_human":"${latency_human}",
        "bytes_in":${bytes_in},
        "bytes_out":${bytes_out}}</code> (default)
      </td>
      <td align="left">
        Log format which can be constructed using the following tags:
        <ul>
          <li>time_unix</li>
		      <li>time_unix_nano</li>
          <li>time_rfc3339</li>
		      <li>time_rfc3339_nano</li>
          <li>id (Request ID - Not implemented)</li>
          <li>remote_ip</li>
          <li>uri</li>
          <li>host</li>
          <li>method</li>
          <li>path</li>
          <li>referer</li>
          <li>user_agent</li>
          <li>status</li>
          <li>latency (In microseconds)</li>
          <li>latency_human (Human readable)</li>
          <li>bytes_in (Bytes received)</li>
          <li>bytes_out (Bytes sent)</li>
          <li>header:&lt;NAME&gt;</li>
          <li>query:&lt;NAME&gt;</li>
          <li>form:&lt;NAME&gt;</li>
        </ul>
      </td>
    </tr>
  </tbody>
</table>

*Example*

`${remote_ip} ${status}`
