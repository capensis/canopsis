import { widgetAlarmsSocketMixin } from '@/mixins/widget/alarms-socket';

describe('widget-alarms-socket', () => {
  describe('methods', () => {
    test('Should not join alarms socket room without alarms', () => {
      const context = {
        $socket: {
          join: jest.fn(),
        },
        alarmsSocketRoom: 'alarms/widget-id',
        updateAlarmInStore: jest.fn(),
      };

      widgetAlarmsSocketMixin.methods.joinToAlarmsSocketRoom.call(context, []);

      expect(context.$socket.join).not.toHaveBeenCalled();
    });

    test('Should join alarms socket room with alarm ids', () => {
      const addListener = jest.fn();
      const context = {
        $socket: {
          join: jest.fn().mockReturnValue({ addListener }),
        },
        alarmsSocketRoom: 'alarms/widget-id',
        updateAlarmInStore: jest.fn(),
      };
      const alarms = [{ _id: 'alarm-1' }, { _id: 'alarm-2' }];

      widgetAlarmsSocketMixin.methods.joinToAlarmsSocketRoom.call(context, alarms);

      expect(context.$socket.join).toHaveBeenCalledWith(
        context.alarmsSocketRoom,
        ['alarm-1', 'alarm-2'],
      );
      expect(addListener).toHaveBeenCalledWith(context.updateAlarmInStore);
    });
  });
});
