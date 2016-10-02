+++
title = "Logger"
[menu.side]
  parent = "plugins"
  weight = 3
+++

## Logger Plugin

### `logger`

Logs HTTP requests

### Configuration

<table>
  <thead>
    <tr>
      <th align="left">Name</th>
      <th align="left">Type</th>
      <th align="left">Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="left"><code>format</code></td>
      <td align="left">string</td>
      <td align="left">
        Log format which can be constructed using the following tags:
        <ul>
          <li>time_rfc3339</li>
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
        </ul>
        Example <code>${remote_ip} ${status}</code>
      </td>
    </tr>
  </tbody>
</table>
