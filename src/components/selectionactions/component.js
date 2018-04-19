Ember.Application.initializer({
    name: 'component-selectionactions',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the selectionactions component for the widget listalarm
         *
         * @class selectionactions component
         */
        var component = Ember.Component.extend({
            tagName: 'td',
            classNames: ['action-cell'],

            /**
             * @property actionsMap
             */
            actionsMap: Ember.A([
                {
                    class: '',
                    mixin_name: 'bulk_pbehavior',
                    caption: 'Apply PBehavior',
					rightName: "listalarm_pbehavior"
                },
                {
                    class: 'glyphicon glyphicon-ok',
                    mixin_name: 'fastack',
                    caption: 'Fast Ack',
					rightName: "listalarm_fastAck"
                },
                {
                    class: 'glyphicon glyphicon-saved',
                    mixin_name: 'ack',
                    caption: 'Ack',
					rightName: "listalarm_ack"
                },
                {
                    class: 'glyphicon glyphicon-ban-circle',
                    mixin_name: 'ackremove',
                    caption: 'Cancel ack',
					rightName: "listalarm_cancelAck"
                },
                {
                    class: 'glyphicon glyphicon-share-alt',
                    mixin_name: 'recovery',
                    caption: 'Restore alarm',
					rightName: 'listalarm_restoreAlarm',
                }
            ]),

			canAction: function(rights, actionName){
				console.error("Rights", rights)
				console.error("ActionName", actionName)
				if (rights.hasOwnProperty(actionName)) {
					if (rights.get(actionName).checksum) {
						return true
					}
				}
				return false
			},


			genAvailableAction: function() {
				var actions = new Array()
				for(i = 0; i < this.get("actionsMap").length; i++) {
					console.error("rigthName", this.actionsMap[i]["rightName"])
					console.error("Can he do this", this.get("canAction")(this.get("rights"), this.actionsMap[i]["rightName"]))

					if (this.get("canAction")(this.get("rights"), this.actionsMap[i]["rightName"])) {
						actions.push(this.actionsMap[i])
					}
				}
				console.error("ACTIONS", actions)
				this.set("availableAction", actions)
			},

			rights: function() {
				return this.get("_parentView._controller.login.rights")
			}.property("rights"),

            /**
             * @method init
             */
            init: function() {
                this._super();
				this.set("rights", this.get("_parentView._controller.login.rights"))
				console.error("Rights ", this.get("rights"))
				this.genAvailableAction()
            },

            actions: {
                /**
                 * @property sendAction
                 */
                sendAction: function (action) {
                    this.sendAction('action', action);
                }
            }

        });

        application.register('component:component-selectionactions', component);
    }
});
