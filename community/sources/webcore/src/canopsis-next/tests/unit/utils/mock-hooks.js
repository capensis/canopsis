/**
 * Mock for date now. Clear yourself after all tests.
 *
 * @param {number} nowTimestamp
 * @deprecated Should be used jest.useFakeTimers({ now: nowTimestamp })
 */
export const mockDateNow = (nowTimestamp) => {
  let dateNowSpy;

  beforeAll(() => {
    dateNowSpy = jest.spyOn(Date, 'now').mockReturnValue(nowTimestamp);
  });

  afterAll(() => {
    dateNowSpy.mockRestore();
  });
};

/**
 * Mock for requestAnimationFrame. Clear yourself after all tests.
 */
export const mockRequestAnimationFrame = () => {
  let requestAnimationFrameSpy = null;

  beforeEach(() => {
    requestAnimationFrameSpy = jest.spyOn(window, 'requestAnimationFrame')
      .mockImplementation(() => {});
  });

  afterEach(() => {
    requestAnimationFrameSpy.mockRestore();
  });
};

/**
 * Mock for date. Clear yourself after all tests.
 *
 * @param {number | Date} nowTimestamp
 */
export const mockDateGetTime = (nowTimestamp) => {
  let dateSpy;

  beforeAll(() => {
    dateSpy = jest
      .spyOn(Date.prototype, 'getTime')
      .mockReturnValue(nowTimestamp);
  });

  afterAll(() => {
    dateSpy.mockRestore();
  });
};

/**
 * Mock for the modals. Clear yourself after all tests.
 */
export const mockModals = () => {
  const modals = {
    show: jest.fn(),
    hide: jest.fn(),
    minimize: jest.fn(),
    maximize: jest.fn(),
    moduleName: 'modals',
  };

  afterEach(() => {
    modals.show.mockReset();
    modals.hide.mockReset();
    modals.minimize.mockReset();
    modals.maximize.mockReset();
  });

  return modals;
};

/**
 * Mock for the popups. Clear yourself after all tests.
 */
export const mockPopups = () => {
  const popups = {
    error: jest.fn(),
    success: jest.fn(),
  };

  afterEach(() => {
    popups.error.mockReset();
    popups.success.mockReset();
  });

  return popups;
};

/**
 * Mock for the router. Clear yourself after all tests.
 */
export const mockRouter = () => {
  const router = {
    push: jest.fn(),
  };

  afterEach(() => {
    router.push.mockReset();
  });

  return router;
};

/**
 * Mock for the sidebar. Clear yourself after all tests.
 */
export const mockSidebar = () => {
  const sidebar = {
    show: jest.fn(),
    hide: jest.fn(),
    moduleName: 'sidebar',
  };

  afterEach(() => {
    sidebar.show.mockReset();
    sidebar.hide.mockReset();
  });

  return sidebar;
};

/**
 * Mock for the socket. Clear yourself after all tests.
 */
export const mockSocket = () => {
  const room = {
    addListener: jest.fn(),
    removeListener: jest.fn(),
  };
  const socket = {
    join: jest.fn().mockReturnValue(room),
    leave: jest.fn().mockReturnValue(room),
  };

  afterEach(() => {
    socket.join.mockClear();
    socket.leave.mockClear();
  });

  return socket;
};

/**
 * Mock for XMLHttpRequest. Clear yourself after all tests.
 */
export const mockXMLHttpRequest = () => {
  const request = {
    send: jest.fn(),
    open: jest.fn(),
    status: undefined,
    responseText: undefined,
    upload: {
      addEventListener: jest.fn(),
    },
  };
  const xmlHttpRequestSpy = jest.spyOn(global, 'XMLHttpRequest').mockReturnValue(request);

  afterEach(() => {
    request.send.mockReset();
    request.open.mockReset();
    request.upload.addEventListener.mockReset();
    request.status = undefined;
    request.responseText = undefined;
  });
  afterAll(() => {
    xmlHttpRequestSpy.mockRestore();
  });

  return request;
};

/**
 * Mock console methods
 *
 * @return {Object}
 */
export const mockConsole = () => {
  const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();

  afterEach(() => {
    consoleErrorSpy.mockReset();
  });

  return {
    error: consoleErrorSpy,
  };
};
