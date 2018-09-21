import { normalize } from 'normalizr';

import queryMixin from '@/mixins/query';
import sideBarMixins from '@/mixins/side-bar/side-bar';
import entitiesViewRowMixin from '@/mixins/entities/view/row';
import entitiesViewWidgetMixin from '@/mixins/entities/view/widget';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

import { viewSchema } from '@/store/schemas';

import { convertUserPreferenceToQuery, convertWidgetToQuery } from '@/helpers/query';
import get from 'lodash/get';

export default {
  props: {
    config: {
      type: Object,
      required: true,
    },
  },
  mixins: [
    queryMixin,
    sideBarMixins,
    entitiesUserPreferenceMixin,
    entitiesViewRowMixin,
    entitiesViewWidgetMixin,
  ],
  data() {
    return {
      clonedView: {

      },
    };
  },
  computed: {
    widget() {
      return this.config.widget;
    },
  },
  methods: {
    isFormValid() {
      if (this.$validator) {
        return this.$validator.validateAll();
      }

      return true;
    },

    prepareSettingsWidget() {
      return this.settings.widget;
    },

    async submit() {
      const isFormValid = await this.isFormValid();

      if (isFormValid) {
        const widget = {
          ...this.widget,
          ...this.prepareSettingsWidget(),
        };

        const userPreference = {
          ...this.userPreference,
          widget_preferences: {
            ...this.userPreference.widget_preferences,
            ...this.settings.widget_preferences,
          },
        };

        const oldRowId = get(widget, '_embedded.parentId');
        const newRowId = this.rowId;

        if (this.rowForCreation && this.rowForCreation._id === this.rowId) {
          this.view.rows.push(this.rowForCreation);
        }

        normalize(this.view, viewSchema); // TODO: DO IT

        this.view.rows = this.view.rows.map((row) => {
          if (newRowId !== oldRowId) {
            if (oldRowId && oldRowId === row._id) {
              return { ...row, widgets: row.widgets.filter(({ _id }) => _id !== widget._id) };
            }

            if (newRowId === row._id) {
              return { ...row, widgets: [...row.widgets, widget] };
            }
          } else if (newRowId === row._id) {
            return {
              ...row,
              widgets: row.widgets.map((oldWidget) => {
                if (oldWidget._id === widget._id) {
                  return widget;
                }

                return oldWidget;
              }),
            };
          }

          return row;
        });

        if (this.rowForCreation && this.rowForCreation === this.rowId) {
          this.createRowInStore({
            row: this.rowForCreation,
          });
        }

        const actions = [this.createUserPreference({ userPreference })];

        if (this.config.isNew) {
          actions.push(this.createWidget({ widget }));
        } else {
          actions.push(this.updateWidget({ widget }));
        }

        await Promise.all(actions);

        this.mergeQuery({
          id: widget._id,
          query: {
            ...convertWidgetToQuery(widget),
            ...convertUserPreferenceToQuery(userPreference),
          },
        });

        this.hideSideBar();
      }
    },
  },
};
