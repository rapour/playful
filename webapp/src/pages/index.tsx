import {useEffect, useState} from 'react'
import Head from 'next/head'
import { Inter } from '@next/font/google'
import styles from '@/styles/Home.module.css'

const inter = Inter({ subsets: ['latin'] })

export default function Home() {

  const [time, setTime] = useState(null)

  useEffect(() => {
    let eventSource = new EventSource(`/api/loc`);
    eventSource.onopen = (e) => { console.log('listen to api-sse endpoint', e)};

    eventSource.onmessage = (e) => {
      //const location = JSON.parse(e.data);
      setTime(e.data)
    };

    eventSource.onerror = (e) => { console.log('error', e )};

    // returned function will be called on component unmount
    return () => {
      eventSource.close();
      // eventSource = null;
    }
  },[])


  return (
    <>
      <Head>
        <title>Playful Particles</title>
        <meta name="description" content="A simple web application to demonstrate hidden markov" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main className={styles.main}>
        <h3 className={inter.className}>
          {time}
        </h3>
      </main>
    </>
  )
}
