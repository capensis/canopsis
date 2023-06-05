export const createModalWrapperStub = className => ({
  template: `
    <div class="${className}">
      <div class="title">
        <slot name="title" />
      </div>
      <div class="text">
        <slot name="text" />
      </div>
      <div class="actions">
        <slot name="actions" />
      </div>
    </div>
  `,
});
