import ws from "k6/ws";
import { check } from "k6";
import { Trend } from "k6/metrics";

let latency = new Trend("latency");

export const options = {
  vus: 2000, // jumlah virtual user
  duration: "30s", // durasi test
};

export default function () {
  const url = "ws://localhost:8081/ws";

  const res = ws.connect(url, {}, function (socket) {
    let start = 0; // deklarasi di luar event supaya bisa diakses semua handler

    socket.on("open", function () {
      start = Date.now();
      socket.send("ping");
    });

    socket.on("message", function (message) {
      let end = Date.now();
      latency.add(end - start);
      // Bisa log kalau mau lihat sample:
      // console.log(`Latency: ${end - start}ms, msg: ${message}`);
    });

    socket.on("close", function () {
      // koneksi tertutup
    });

    socket.setTimeout(function () {
      socket.close();
    }, 5000);
  });

  check(res, { "status is 101": (r) => r && r.status === 101 });
}
