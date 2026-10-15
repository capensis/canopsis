export const createFormStub = className => ({
  template: `
    <form class="${className}" v-on="$listeners">
      <slot />
    </form>
  `,
});

export const getFormBlockArrayFieldStub = () => ({
  props: {
    form: {
      type: Array,
      default: () => [],
    },
    value: {
      type: Array,
      default: () => [],
    },
    itemToForm: {
      type: Function,
      required: true,
    },
  },
  computed: {
    items() {
      return this.form?.length ? this.form : (this.value ?? []);
    },
  },
  template: `
    <div>
      <div v-for="(item, index) in items" :key="index">
        <slot name="item" :index="index" :remove="() => remove(index)" />
      </div>
      <button type="button" @click="add">add</button>
    </div>
  `,
  methods: {
    add() {
      this.$emit('input', [...this.items, this.itemToForm()]);
    },
    remove(index) {
      this.$emit('input', this.items.filter((_, itemIndex) => itemIndex !== index));
    },
  },
});

export const getFormGeneralPatternsTabsStub = () => ({
  template: `
    <div>
      <slot name="general" :set-ref="noop" />
      <slot name="patterns" :set-ref="noop" />
      <slot name="additional" :set-ref="noop" />
      <slot name="test-query" />
    </div>
  `,
  methods: {
    noop() {},
  },
});
