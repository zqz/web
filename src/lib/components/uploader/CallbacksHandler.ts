function CallbackHandler<T>() {
  let callbacks = new Map<T, Function>();

  function on(eventName: T, func: Function) {
    callbacks.set(eventName, func);
  }

  function emit(eventName: T, ...args: any[]) {
    const cb = callbacks.get(eventName);

    if (cb === undefined || cb === null) {
      return;
    }

    cb(...args);
  }

  function call(eventName: T, ...args: any[]) {
    const cb = callbacks.get(eventName);

    if (cb === undefined || cb === null) {
      return;
    }

    cb(...args);
  }

  return { on, emit, call };
}

export default CallbackHandler;
