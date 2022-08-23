import { MODALS } from '@/constants';

import { flowchartPointToForm } from '@/helpers/forms/map';

export const pointsMixin = {
  props: {
    points: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      shownMenu: false,
      shownPointDialog: false,
      clientX: 0,
      clientY: 0,
      addingPoint: undefined,
      editingPoint: undefined,
    };
  },
  computed: {
    isDialogOpened() {
      return this.shownPointDialog;
    },

    contextmenuItems() {
      if (this.editingPoint) {
        return [
          {
            text: this.$t('map.editPoint'),
            action: this.openEditPointDialog,
          },
          {
            text: this.$t('map.removePoint'),
            action: this.showRemovePointModal,
          },
        ];
      }

      return [
        {
          text: this.$t('map.addPoint'),
          action: this.openAddPointDialog,
        },
      ];
    },
  },
  methods: {
    emitUpdatePoints(data) {
      this.$emit('update:points', data);
    },

    updatePointInModel(data) {
      this.emitUpdatePoints(
        this.points.map(point => (point._id === data._id ? data : point)),
      );
    },

    addPointToModel(data) {
      this.emitUpdatePoints([...this.points, data]);
    },

    removePointFromModel(data) {
      this.emitUpdatePoints(
        this.points.filter(point => data._id !== point._id),
      );
    },

    openContextmenu() {
      this.shownMenu = true;
    },

    closeContextmenu() {
      this.shownMenu = false;
    },

    openPointDialog() {
      this.shownPointDialog = true;
    },

    closePointDialog() {
      this.shownPointDialog = false;
      this.addingPoint = undefined;
      this.editingPoint = undefined;
    },

    setOffsetsByEvent(event) {
      this.clientX = event.clientX;
      this.clientY = event.clientY;
    },

    setOffsetsByPoint(point) {
      const { x, y } = this.$refs.svg.getBoundingClientRect();

      /**
       * TODO: Should be fixed scale
       * @type {number}
       */
      this.clientX = x + point.x - this.viewBoxObject.x;
      this.clientY = y + point.y - this.viewBoxObject.y;
    },

    handleContextmenu(event) {
      if (this.shownMenu || this.isDialogOpened) {
        return;
      }

      this.setOffsetsByEvent(event);
      this.openContextmenu();
    },

    handleEditContextmenu(point) {
      if (this.shownMenu || this.isDialogOpened) {
        return;
      }

      this.setOffsetsByPoint(point);
      this.openContextmenu();

      this.editingPoint = point;
    },

    openAddPointDialog() {
      const { x, y } = this.normalizeCursor({ x: this.clientX, y: this.clientY });

      this.addingPoint = flowchartPointToForm({
        x,
        y,
      });

      this.closeContextmenu();
      this.openPointDialog();
    },

    openEditPointDialog() {
      this.closeContextmenu();
      this.openPointDialog();
    },

    submitPointDialog(data) {
      if (this.editingPoint) {
        this.updatePointInModel(data);
      } else {
        this.addPointToModel(data);
      }

      this.closePointDialog();
    },

    showRemovePointModal() {
      this.closeContextmenu();

      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => {
            this.removePointFromModel(this.editingPoint);

            this.closePointDialog();
          },
        },
      });
    },
  },
};
