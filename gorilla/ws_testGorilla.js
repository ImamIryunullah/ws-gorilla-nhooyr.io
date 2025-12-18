import ws from "k6/ws";
import { check } from "k6";
import { Trend, Counter, Rate } from "k6/metrics";

// Existing metrics
let latency = new Trend("latency", true); // true = collect pXX
let messages_sent = new Counter("messages_sent");
let messages_received = new Counter("messages_received");
let connection_errors = new Counter("connection_errors");
let message_errors = new Counter("message_errors");
let success_rate = new Rate("success_rate");

// 🔥 New metrics
let connect_time = new Trend("connect_time");
let disconnect_time = new Trend("disconnect_time");
let ping_rtt = new Trend("ping_rtt");
let vu_active = new Counter("vu_active");
let vu_completed = new Counter("vu_completed");
let unexpected_close = new Counter("unexpected_close");
let message_size = new Trend("message_size");

// Custom metrics for throughput
let throughput = new Trend("throughput");

export const options = {
  vus: 3000,
  duration: "60s",
  thresholds: {
    // Tambahkan threshold untuk P99 latency
    latency: ["p(99)<500"], // P99 latency harus <500ms
    messages_received: ["count>0"],
  },
};

export default function () {
  const url = "wss://websocket.pollacheialnetworks.my.id/gorilla";
  const connectStart = Date.now();

  const res = ws.connect(url, {}, function (socket) {
    vu_active.add(1);

    let start = 0;

    socket.on("open", () => {
      connect_time.add(Date.now() - connectStart);

      start = Date.now();
      socket.send("ping");
      messages_sent.add(1);
    });

    socket.on("message", (message) => {
      let end = Date.now();
      let rtt = end - start;

      latency.add(rtt);
      ping_rtt.add(rtt);

      if (!message) {
        message_errors.add(1);
      } else {
        messages_received.add(1);
        message_size.add(message.length);
      }

      // Hitung throughput sebagai pesan per detik
      throughput.add(
        messages_received.value / ((Date.now() - connectStart) / 1000)
      );

      success_rate.add(1);

      vu_completed.add(1);
      socket.close();
    });

    socket.on("error", (e) => {
      connection_errors.add(1);
      success_rate.add(0);
    });

    socket.on("close", () => {
      disconnect_time.add(Date.now() - start);

      if (messages_received.value === 0) {
        unexpected_close.add(1);
      }
    });

    // fallback force close
    socket.setTimeout(() => socket.close(), 5000);
  });

  check(res, { "status is 101": (r) => r && r.status === 101 }) ||
    connection_errors.add(1);
}
