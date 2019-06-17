<template lang="pug">
  div
    v-layout(justify-end)
      v-tooltip(v-for="(action, index) in mainActions", :key="`action-${index}`", top)
        v-btn(slot="activator", icon, @click="action.action()")
          v-icon(:class="action.iconClass") {{ action.icon }}
        span {{ action.tooltip }}
    v-layout(row)
      v-flex(xs12)
        v-treeview(:items="treeviewItems", :open.sync="opened")
          template(slot="label", slot-scope="{ item }")
            v-flex(xs12)
              v-layout(row)
                v-flex(xs6) {{ item.name }}
                  span(v-show="item.value") :
                template(v-if="item.value")
                  v-flex(v-if="isSimpleRule(item.value)")
                    span.body-1.font-italic {{ item.value }}
                  v-flex(v-else)
                    v-layout(column)
                      v-flex(v-for="(field, fieldKey) in item.value", :key="fieldKey")
                        p.body-1.font-italic {{ fieldKey }} {{ field }}
          template(slot="append", slot-scope="{ item }")
            div
              v-tooltip(v-for="(action, index) in getActionsForItem(item)", :key="`action-${index}`", top)
                v-btn(slot="activator", icon, @click="action.action(item)")
                  v-icon(:class="action.iconClass") {{ action.icon }}
                span {{ action.tooltip }}
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
      actionsMap: {
        addRuleField: {
          tooltip: 'Add rule field',
          icon: 'add',
          iconClass: 'primary--text',
          action: this.showAddRuleFieldModal,
        },
        editRuleField: {
          tooltip: 'Edit rule field',
          icon: 'edit',
          action: this.showEditRuleModal,
        },
        addObjectRuleField: {
          tooltip: 'Add object rule field',
          icon: 'library_add',
          iconClass: 'primary--text',
          action: this.showAddObjectFieldModal,
        },
        editObjectRuleField: {
          tooltip: 'Edit object rule field',
          icon: 'edit',
          action: this.showEditObjectFieldModal,
        },
        removeRuleField: {
          tooltip: 'Remove rule field',
          icon: 'remove',
          iconClass: 'error--text',
          action: this.deleteRule,
        },
      },
    };
  },
  computed: {
    mainActions() {
      const { actionsMap } = this;

      return [
        actionsMap.addRuleField,
        actionsMap.addObjectRuleField,
      ];
    },

    getActionsForItem() {
      const { actionsMap } = this;

      return (item) => {
        if (item.value) {
          return [
            actionsMap.editRuleField,
            actionsMap.removeRuleField,
          ];
        }

        return [
          actionsMap.addRuleField,
          actionsMap.addObjectRuleField,
          actionsMap.editObjectRuleField,
          actionsMap.removeRuleField,
        ];
      };
    },

    treeviewItems() {
      return this.parseObjectToTreeview(this.pattern);
    },

    isSimpleRule() {
      return rule => !isObject(rule);
    },
  },
  methods: {
    isValueRule(rule) {
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
        };

        if (this.isValueRule(value)) {
          item.value = value;
        } else {
          item.children = this.parseObjectToTreeview(value, path);
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
          field: {
            label: 'Field',
            validationRules: 'required',
            name: 'field',
          },
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
          field: {
            label: 'Field',
            value: item.name,
            validationRules: 'required',
            name: 'field',
          },
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
