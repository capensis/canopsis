export default {
  short: {
    year: 'numeric', month: 'short', day: 'numeric',
  },
  long: {
    year: 'numeric',
    month: '2-digit',
    hour12: true,
    pattern: '{day}-{month}-{year} {hour}:{minute}:{second}',
    pattern12: '{hour}:{minute}:{second} {ampm}',
  },
  time: {
    hour: 'numeric', minute: 'numeric', second: 'numeric',
  },
};
