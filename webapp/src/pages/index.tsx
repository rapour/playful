import { useEffect, useState } from "react";
import Head from "next/head";
import { Inter } from "@next/font/google";
import styles from "@/styles/Home.module.css";

const inter = Inter({ subsets: ["latin"] });

const colors = ["green", "white", "red"]

const getColor = (id: any) => colors[id % 3]

export default function Home() {
  const [nodes, setNodes] = useState([]);

  useEffect(() => {
    let eventSource = new EventSource(`/api/loc`);
    eventSource.onopen = (e) => {
      console.log("listen to api-sse endpoint", e);
    };

    eventSource.onmessage = (e) => {
      const locations = JSON.parse(e.data);
      let ns = locations.map((e: any) => {
        let n: any = nodes.find((a: any) => a.id === e.id)
        if (e.visible || n == undefined) {
          return {id: e.id, top: e.alt, left: e.lon, visible: e.visible}
        } else {
          n.visible = false
          return n
        }
      })
      
      setNodes(ns);
    };

    eventSource.onerror = (e) => {
      console.log("error", e);
    };

    // returned function will be called on component unmount
    return () => {
      eventSource.close();
      // eventSource = null;
    };
  }, []);

  return (
    <>
      <Head>
        <title>Playful Particles</title>
        <meta
          name="description"
          content="A simple web application to demonstrate hidden markov"
        />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main className={styles.main}>
        <div className={styles.field}>
          {nodes.map((elm: any) => (
            <div
              key={elm.id}
              className={styles.dot}
              style={{
                top: elm.top,
                left: elm.left,
                display: elm.visible ? "block" : "none",
                backgroundColor: getColor(elm.id)
              }}
            />
          ))}
        </div>
      </main>
    </>
  );
}
