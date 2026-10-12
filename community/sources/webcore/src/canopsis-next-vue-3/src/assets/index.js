const modules = import.meta.glob('@/assets/event-filter-actions-types/*.svg', {
  eager: true,
  import: 'default',
});

export const eventFilterActionsTypesImages = Object.fromEntries(
  Object.entries(modules).map(([path, url]) => {
    const fileName = path
      .split('/')
      .pop();

    return [`./${fileName}`, url];
  }),
);
