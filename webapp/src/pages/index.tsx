import { useEffect, useState } from "react";
import {motion} from 'framer-motion'
import Head from "next/head";
import { Inter } from "@next/font/google";
import styles from "@/styles/Home.module.css";

const inter = Inter({ subsets: ["latin"] });

export default function Home() {
  const [top, setTop] = useState(100);
  const [left, setLeft] = useState(100);
  

  useEffect(() => {
    let eventSource = new EventSource(`/api/loc`);
    eventSource.onopen = (e) => {
      console.log("listen to api-sse endpoint", e);
    };

    eventSource.onmessage = (e) => {
      const location = JSON.parse(e.data);
      setTop(2 * location.Altitude)
      setLeft(2 * location.Longitude)
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
        <motion.div
        className={styles.dot}
            animate={{ top: top, left: left }}
            transition={{ delay: 0.5 }}
          />
          {/* <div className={styles.dot} style={{ top: top, left: left }} /> */}
        </div>
      </main>
    </>
  );
}
