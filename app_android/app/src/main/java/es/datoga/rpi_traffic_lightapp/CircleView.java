package es.datoga.rpi_traffic_lightapp;

import android.content.Context;
import android.content.res.TypedArray;
import android.graphics.Canvas;
import android.graphics.Color;
import android.graphics.Paint;
import android.util.AttributeSet;
import android.util.Log;
import android.view.MotionEvent;
import android.view.View;

public class CircleView extends View {
    private static final String TAG = CircleView.class.getSimpleName();

    public static int DEFAULT_DIMENSION = 20;

    private int onColor;
    private float dimension;

    private int mWidth;
    private int mHeight;

    Paint onCirclePaint;
    Paint offCirclePaint;
    Paint activePaint;


    public CircleView(Context context) {
        super(context);
        Log.d(TAG, "Create circle view");
        init(null, 0);
    }

    public CircleView(Context context, AttributeSet attrs) {
        super(context, attrs);
        Log.d(TAG, "Create circle view");
        init(attrs, 0);
    }

    public CircleView(Context context, AttributeSet attrs, int defStyle) {
        super(context, attrs, defStyle);
        Log.d(TAG, "Create circle view");
        init(attrs, defStyle);
    }

    private void init(AttributeSet attrs, int defStyle) {
        Log.d(TAG, "Init circle view");

        // Load attributes
        final TypedArray a = getContext().obtainStyledAttributes(
                attrs, R.styleable.CircleView, defStyle, 0);

        onColor = a.getColor(R.styleable.CircleView_onColor, Color.GRAY);

        a.recycle();

        onCirclePaint = new Paint(Paint.ANTI_ALIAS_FLAG);
        onCirclePaint.setColor(onColor);

        Log.d(TAG, "On color " + onColor);

        offCirclePaint = new Paint(Paint.ANTI_ALIAS_FLAG);
        offCirclePaint.setColor(Color.GRAY);

        activePaint = offCirclePaint;
    }

    @Override
    protected void onDraw(Canvas canvas) {
        super.onDraw(canvas);

        canvas.drawCircle(getMeasuredWidth()/2, getMeasuredHeight()/2, getMeasuredHeight()/2, this.activePaint);
    }

    @Override protected void onMeasure(int widthMeasureSpec, int heightMeasureSpec)
    {
        mWidth = View.MeasureSpec.getSize(widthMeasureSpec);
        mHeight = View.MeasureSpec.getSize(heightMeasureSpec);

        setMeasuredDimension(mWidth, mHeight);
    }

    public void switchOn() {
        this.activePaint = this.onCirclePaint;
        invalidate();
    }

    public void switchOff() {
        this.activePaint = this.offCirclePaint;
        invalidate();
    }
}
