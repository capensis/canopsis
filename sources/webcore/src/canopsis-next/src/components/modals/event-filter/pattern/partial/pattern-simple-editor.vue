<template lang="pug">
  div
    v-layout(justify-end)
      v-tooltip(top)
        v-btn(slot="activator", icon, @click="showAddRuleFieldModal()")
          v-icon.primary--text add
        span Add field with value
      v-tooltip(top)
        v-btn(slot="activator", icon, @click="showAddObjectFieldModal()")
          v-icon.primary--text library_add
        span Add object field
    v-layout(row)
      v-flex(xs12)
        v-treeview.position-relative(:items="treeviewItems", :open.sync="opened")
          template(slot="label", slot-scope="{ item }")
            v-flex.field-treeview__label(xs12)
              v-layout(row)
                v-flex(xs6) {{ item.name }}
                v-flex(v-if="item.simple", xs6)
                  span(v-if="isValueObject(item.value)")
                    div(v-for="(rule, ruleKey) in item.value", :key="ruleKey")
                      span {{ ruleKey }}: {{ rule }}
                  span(v-else) {{ item.value }}
          template(slot="append", slot-scope="{ item }")
            .field-treeview__actions
              template(v-if="!item.simple")
                v-btn(icon, small, @click="showAddRuleFieldModal(item)")
                  v-icon.primary--text add
                v-btn(icon, small, @click="showAddObjectFieldModal(item)")
                  v-icon.primary--text library_add
              v-btn(v-if="item.simple", icon, small, @click="showEditRuleModal(item)")
                v-icon edit
              v-btn(v-else, icon, small, @click="showEditObjectFieldModal(item)")
                v-icon edit
              v-btn(icon, small, @click="deleteRule(item)")
                v-icon.error--text remove
</template>

<script>
import { isObject } from 'lodash';

import { MODALS } from '@/constants';

import formMixin from '@/mixins/form';
import modalMixin from '@/mixins/modal';

export default {
  mixins: [formMixin, modalMixin],
  model: {
    prop: 'pattern',
    event: 'input',
  },
  props: {
    pattern: {
      type: Object,
      required: true,
    },
    operators: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      opened: [],
    };
  },
  computed: {
    treeviewItems() {
      return this.parseObjectToTreeview(this.pattern);
    },

    isValueObject() {
      return value => isObject(value);
    },
  },
  methods: {
    isSimpleRule(rule) {
      if (isObject(rule)) {
        const items = Object.entries(rule);

        return items.length && items.every(([key, value]) => this.operators.indexOf(key) !== -1 && !isObject(value));
      }

      return true;
    },

    parseObjectToTreeview(object, prevPath = []) {
      return Object.entries(object).map(([key, value]) => {
        const path = [...prevPath, key];
        const item = {
          path,
          name: key,
          id: path.join('.'),
          simple: this.isSimpleRule(value),
        };

        if (!item.simple) {
          item.children = this.parseObjectToTreeview(value, path);
        } else {
          item.value = value;
        }

        return item;
      }, []);
    },

    deleteRule(item) {
      this.removeField(item.path);
    },

    showEditRuleModal(item) {
      this.showModal({
        name: MODALS.addEventFilterRuleToPattern,
        config: {
          ruleKey: item.name,
          ruleValue: item.value,
          isSimpleRule: !isObject(item.value),
          operators: this.operators,
          action: newRule => this.updateAndMoveField([item.name], [newRule.field], newRule.value),
        },
      });
    },

    showAddObjectFieldModal(parent) {
      const parentPath = parent ? parent.path : [];

      this.showModal({
        name: MODALS.textFieldEditor,
        config: {
          title: 'Add object field',
          validationRules: 'required',
          action: (field) => {
            this.updateField([...parentPath, field], {});

            this.$nextTick(() => this.openTreeviewItem(parent));
          },
        },
      });
    },

    showEditObjectFieldModal(item) {
      this.showModal({
        name: MODALS.textFieldEditor,
        config: {
          title: 'Add object field',
          validationRules: 'required',
          action: field => this.moveField([...item.path], [field]),
        },
      });
    },

    showAddRuleFieldModal() {
      this.showModal({
        name: MODALS.addEventFilterRuleToPattern,
        config: {
          operators: this.operators,
          action: newRule => this.updateField([newRule.field], newRule.value),
        },
      });
    },

    openTreeviewItem(item) {
      if (item && this.opened.indexOf(item.id) === -1) {
        this.opened.push(item.id);
      }
    },
  },
};
</script>
