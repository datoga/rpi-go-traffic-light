package es.datoga.rpi_traffic_lightapp;

import android.app.Activity;
import android.content.Context;
import android.content.res.Resources;
import android.os.Bundle;
import android.os.Handler;

import android.util.Log;
import android.view.View;
import android.widget.CompoundButton;
import android.widget.Switch;

import actuators.Actuators;
import actuators.RemoteActuator;

public class RPICanvas extends Activity {
    private static final String TAG = RPICanvas.class.getSimpleName();

    CircleView greenCircle;
    CircleView redCircle;
    CircleView yellowCircle;
    Switch manualAutoSwitch;

    Handler customHandler;

    RemoteActuator remoteActuator;

    Runnable initToAutoThread = new Runnable() {
        @Override
        public void run() {
            try {
                remoteActuator.setAutomatic();
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
    };

    Runnable updateTimerThread = new Runnable()
    {
        @Override
        public void run() {

            Log.d(TAG, "Executing handler");

            try {
                String remoteState = remoteActuator.getRemoteState();

                Log.d(TAG, "Remote state: " + remoteState);

                switch (remoteState) {
                    case "green": setGreen(); break;
                    case "red": setRed(); break;
                    case "yellow": setYellow(); break;

                    default: Log.e(TAG, "Unhandled remote state " + remoteState);
                }

            } catch (Exception e) {
                e.printStackTrace();
                try {
                    Log.w(TAG, "Remote actuator was disconnected, trying to reconnect");
                    remoteActuator.start();
                } catch (Exception e1) {
                    e1.printStackTrace();
                }
            }

            customHandler.postDelayed(this, 1000);
        }
    };

    private View.OnClickListener circleViewListener = new View.OnClickListener() {
        @Override
        public void onClick(View v) {
            Log.d(TAG, "Clicked view " + v.getId() );
            manualAutoSwitch.setChecked(false);

            try {
                if (v.getId() == greenCircle.getId()) {
                    setGreen();
                    remoteActuator.setGreen();
                } else if (v.getId() == yellowCircle.getId()) {
                    setYellow();
                    remoteActuator.setYellow();
                } else if (v.getId() == redCircle.getId()) {
                    setRed();
                    remoteActuator.setRed();
                } else {
                    Log.w(TAG, "View " + v.getId() + " not expected");
                }
            }
            catch (Exception e) {
                e.printStackTrace();
            }
        }
    };

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_rpicanvas);

        greenCircle = (CircleView) findViewById(R.id.greencircle);
        yellowCircle = (CircleView) findViewById(R.id.yellowcircle);
        redCircle = (CircleView) findViewById(R.id.redcircle);

        greenCircle.setOnClickListener(circleViewListener);
        yellowCircle.setOnClickListener(circleViewListener);
        redCircle.setOnClickListener(circleViewListener);

        remoteActuator = Actuators.newRemoteActuator();
        try {
            remoteActuator.start();
        } catch (Exception e) {
            Log.e(TAG, e.toString() );
            e.printStackTrace();
        }

        Handler initToAutoHandler = new Handler();
        initToAutoHandler.postDelayed(initToAutoThread, 2000);

        customHandler = new Handler();
        customHandler.postDelayed(updateTimerThread, 3000);

        manualAutoSwitch = (Switch) findViewById(R.id.manual_auto_switch);

        manualAutoSwitch.setOnCheckedChangeListener(new CompoundButton.OnCheckedChangeListener() {
            public void onCheckedChanged(CompoundButton buttonView, boolean isChecked) {
                try {
                    if (isChecked) {
                        Log.d(TAG, "Setting auto mode");
                        remoteActuator.setAutomatic();
                    } else {
                        Log.d(TAG, "Setting manual mode");
                        remoteActuator.setManual();
                    }
                }
                catch (Exception e) {
                    e.printStackTrace();
                    buttonView.setChecked(!isChecked);
                }
            }
        });
    }

    protected void setGreen() {
        greenCircle.switchOn();
        yellowCircle.switchOff();
        redCircle.switchOff();

    }

    protected void setYellow() {
        greenCircle.switchOff();
        yellowCircle.switchOn();
        redCircle.switchOff();
    }

    protected void setRed() {
        greenCircle.switchOff();
        yellowCircle.switchOff();
        redCircle.switchOn();
    }
}
