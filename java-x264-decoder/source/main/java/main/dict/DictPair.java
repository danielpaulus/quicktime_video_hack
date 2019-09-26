package main.dict;

public class DictPair {

	private final String key;
	private final Object value;

	public DictPair(String key, Object value) {
		this.key = key;
		this.value = value;
	}

	@Override
	public String toString() {
		return "{" + key + ":" + value + "}";
	}
}
