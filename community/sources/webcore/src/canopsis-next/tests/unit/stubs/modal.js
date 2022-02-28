export const createModalWrapperStub = className => ({
  template: `
    <div class="${className}">
      <slot name="title" />
      <slot name="text" />
      <slot name="actions" />
    </div>
  `,
});
