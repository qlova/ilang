import java.util.ArrayList;

class Arrayception {
	public static void main(String[] args) {
		ArrayList a = new ArrayList();
	
		for (int i = 0; i < 1000000; i++) {
			a.add(i*2);
		}
	
		System.out.println(a.get(12317));
	}
}
