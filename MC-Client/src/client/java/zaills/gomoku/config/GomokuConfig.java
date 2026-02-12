package zaills.gomoku.config;

import me.shedaniel.autoconfig.ConfigData;
import me.shedaniel.autoconfig.annotation.Config;

@Config(name = "gomoku")
public class GomokuConfig implements ConfigData {
	public String baseURL = "http://127.0.0.1:8080";
}
