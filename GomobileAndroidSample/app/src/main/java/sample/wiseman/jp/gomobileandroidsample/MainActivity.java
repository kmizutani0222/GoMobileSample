package sample.wiseman.jp.gomobileandroidsample;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

import go.postSample.PostSample;

public class MainActivity extends AppCompatActivity implements View.OnClickListener {

    private EditText param1;
    private EditText param2;
    private TextView response;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        PostSample.setDEBUG_KBN(0);

        response = (TextView) findViewById(R.id.response);
        param1 = (EditText) findViewById(R.id.param1);
        param2 = (EditText) findViewById(R.id.param2);

        Button request_btn = (Button) findViewById(R.id.request_btn);
        request_btn.setOnClickListener(this);
    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            // リクエストボタンが押された
            case R.id.request_btn:
                PostSample.Initialize();
                PostSample.SetParams("param1", param1.getText().toString());
                PostSample.SetParams("param2", param2.getText().toString());
                response.setText(PostSample.HttpPost("/apisample"));
                break;
        }
    }
}
