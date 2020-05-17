const CallbackHandler = () => {
  let callbacks = {};

  function on(callback, func) {
    callbacks[callback] = func;
  }

  function call(name, ...args) {
    const cb = callbacks[name];
    if ( cb === undefined) {
      return;
    }

    cb(...args);
  }

  return { on, call };
}

export default CallbackHandler;
