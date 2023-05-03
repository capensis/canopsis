export const expectsOneInput = (wrapper, data) => {
  const inputEvents = wrapper.emitted('input');

  expect(inputEvents).toHaveLength(1);

  const [eventData] = inputEvents[0];
  expect(eventData).toEqual(data);
};
