(() => {
  const chart = document.getElementById("chart")
  const temps = document.getElementById("temperatures")
  const calibrations = [0,1]
  const strings = ['food', 'method']

  const startStream = () => {
    const socket = new WebSocket(`ws://${window.location.host}/api/stream`)
    socket.onmessage = (event) => {
      const reading = JSON.parse(event.data)
      reading.received = new Date(Date.parse(reading.received))
      chart.src = `/api/chart.png?v=${reading.received.getTime()}`

      temps.innerHTML = ''
      
      reading.temperatures.forEach(temp => {
        const el = document.createElement("div")
        el.className = 'temperature'
        el.innerHTML = `${temp.toFixed(2)}&deg; F`
        temps.appendChild(el)
      })
    }
  }

  const loadMetadata = () => {
    return fetch('/api/metadata')
      .then(res => res.json())
      .then(info => {
        strings.forEach(key => {
          document.getElementById(key).value = info[key]
        })
        calibrations.forEach(idx => {
          document.getElementById(`calibration${idx}`).value = info.calibrations[idx]
        })
      })
  }

  const registerSubmitHandler = () => {
    document.getElementById('form').addEventListener('submit', e => {
      e.preventDefault()
      const md = {
        calibrations: []
      }
      strings.forEach(key => {
        md[key] = document.getElementById(key).value
      })
      calibrations.forEach(idx => {
        md.calibrations[idx] = parseFloat(document.getElementById(`calibration${idx}`).value)
      })
      fetch('/api/metadata', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(md)
      }).catch(e => console.error(e))
      return false
    })
  }

  startStream()
  loadMetadata()
    .then(() => registerSubmitHandler())
    .catch(e => console.error(e))
})()
