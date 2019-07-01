# Java 线程同步

Java 中通过 `implements Runnable` 能够使一个类变成**可运行类**。通过重载 `run` 方法，来表明可运行的操作。

之后通过 `Thread` 类的构造函数，创建一个线程，调用 `start` 方法启动线程。

```java
// MyThread 是自定义的可运行类，第二个参数是线程的名字，在可运行类中通过 Thread.currentThread().getName() 获取到
Thread myThread = new Thread(MyThread, 'ThreadName');
// 执行线程，即运行可运行类中的 run() 方法
myThread.start();
```

此外也可以通过 `extends Thread` 来直接实现线程的子类。这样通过重载 `start` 方法，也可以创建线程。但一般不推荐这么做。

## 线程同步问题

当多个线程访问同一个共享变量时，会发生同步问题。即可能在一个线程读取变量，修改并写回的三个步骤中被其他线程插入，导致结果不可预知。

解决方案有两种，分别是乐观锁和悲观锁。

1. 乐观锁表示可以随意读取，但写入时需要进行一些检查，如果检查通过则继续；不通过就失败，等待重试。这适用于冲突较少的系统中。

    通常乐观锁使用版本号或者 CAS 来实现。Java 中的 `AtomaticNumber` 等类就是用 CAS 实现的乐观锁。CAS 底层使用的是系统命令 CMPXCHG，这是一个原子操作。

2. 悲观锁表示不论读取还是写入，都必须先获取排他的锁。只有一个线程可以获得锁并进行操作，其他线程必须等待。

    悲观锁在 Java 中的实现就是 `synchronized`。

## synchronized

Java 的 `synchronized` 可以锁定两种内容：

1. 锁定对象。通过两种方式：修饰对象 `synchronized(this|object){}` 或者修饰非静态方法。一个类可以有多个实例，非静态方法存在于每个实例上，因此修饰非静态方法也属于锁定在对象上。

2. 锁定类。同样通过两种方式：修饰类对象 `synchronized(MyClass.class){}` 或者修饰静态方法。静态方法存在于类对象中。

### 什么都不锁定

如果线程中什么都不锁定，则如开头所说，执行顺序不定，可以随意被挂起，插入执行其他线程的代码。

也因此，如果是非阻塞的异步任务，就可以不锁定。（或者锁定其中部分关键操作）

### 锁定对象

在创建线程时，如果传入的是同一个可运行类的实例，那么在可运行类中通过锁定 `this` 对象就可以确保这些线程之间互相同步。

```java
// MyRunnable implements Runnable with function 'run()'
MyRunnable myRunnable = new MyRunnable();
// 下面两个线程创建时使用了同一个可运行类的实例。所以在可运行类中的 this 都指向 myRunnable，是同一个。
Thread thread1 = new Thread(myRunnable, "Thread1");
Thread thread2 = new Thread(myRunnable, "Thread2");
```

在 `synchronized(this){}` 代码块内的代码互相排斥，只能有一个线程执行，而且执行结束了才释放锁。而在代码块外部的代码可以并发。

通过在非静态方法上添加 `synchronized` 修饰符可以对方法进行限制，这样这个方法整体就是同步的了。相当于整个方法内的代码包在 `synchronized(this){}` 中。

此外，两种写法之间也是互相同步的，因为他们获取的是同一个锁（可以认为都是 `this`，也就是 `myRunnable`）。

当创建线程时传入的是不同的实例，那么锁定 `this` 就失效了，因为这些 `this` 也是互不相同的，他们拿着不同的锁。

```java
// 下面两个线程创建时使用了不同的可运行类的实例。所以在可运行类中的 this 不是同一个，同步代码块会别穿插执行。
Thread thread1 = new Thread(new MyRunnable(), "Thread1");
Thread thread2 = new Thread(new MyRunnable(), "Thread2");
```

### 锁定类

通过 `synchronized(MyRunnable.class)` 或者 `private synchronized static void foo()` 来对类进行锁定。

因为所有实例都是由同一个类创建的。那么锁定了类的话，不管是同一个实例还是不同的实例，都能够做到互相同步。

### 两种混用

类的锁存在于类对象中，对象锁存在于实例对象中，因此虽然类和实例有关系，但两个锁是独立的。**因此这两种锁混用的时候，不会互相影响**。

## 补充

1. `synchronized`关键字不能继承。

    对于父类中的 synchronized 修饰方法，子类在覆盖该方法时，默认情况下不是同步的，必须显示的使用 synchronized 关键字修饰才行。

2. 在定义接口方法时不能使用 `synchronized` 关键字。

3. 构造方法不能使用 `synchronized` 关键字，但可以使用 `synchronized` 代码块来进行同步。

## 可重入锁

考虑 `synchornized` 的实现方式。如果简单的使用一个 boolean 变量来记录是否取得锁，并且将没有取得锁的线程阻挡的话，可能会发生如下的死锁问题：

```java
public class ReentrantTest implements Runnable {

    public synchronized void get() {
        System.out.println(Thread.currentThread().getName());
        set();
    }

    public synchronized void set() {
        System.out.println(Thread.currentThread().getName());
    }

    public void run() {
        get();
    }

    public static void main(String[] args) {
        ReentrantTest rt = new ReentrantTest();
        for(int i = 0; i < 100; i++) {
            new Thread(rt).start();
        }
    }
}
```

在如上代码中，`get` 方法需要获取类锁，之后执行。但其内部调用的 `set` 方法也需要获取类锁，且和 `get` 获取的那个相同。假设 `synchornized` 的实现真如上面所说的那么简单，那么这里就会产生死锁，因为 `set` 拿不到锁。但实际上，锁又已经被这个线程获取到了（在 `get` 里面），相当于自己等自己。

但实际执行起来，这段代码不会有如上问题。原因在于 `synchornized` 的实现中还会额外记录获取到锁的线程是哪一个，避免这种“自己等自己”的情况出现。这种就叫做**可重入性**，即同一个进程如果获取到了锁，再获取时同一个锁的时候不会被卡住，依然能够继续执行。这样在递归调用时不会发生死锁。

此外，`synchornized` 还会记录被锁的次数。像上述代码中 `get`, `set` 分别调用一次的，内部会有一个 `count` 记录锁被获取的次数是 2。只有当解锁 2 次之后才会真正解锁。

## 其他锁的类型

Java 中锁可以分为三类，分别是：

1. 偏向锁：最轻量级的锁。当一个线程获取偏向锁之后，运行完同步代码后不会立即释放，以保证下次再进入同步代码时不用重新获取锁，可以直接运行。而当竞争出现的时候（其他线程也要获取偏向锁），线程才会释放锁，称为锁撤销。相当于最大限度的让一个线程持有锁，因为它预设绝大部分情况是同一个线程要获取同一个锁。

    优点：加锁和解锁不需要额外小号，和执行非同步方法相比仅存在纳秒级的差距

    缺点：如果线程之间存在竞争，会带来额外的锁撤销的消耗

    适用：（绝大部分情况）只有一个线程访问同步代码

2. 轻量级锁：中等量级的锁，其实就是乐观锁。如果没有获取到锁，则线程自旋，并不退出，省去了线程休眠和唤醒的消耗。CAS就是这一种。

    优点：竞争的线程不会阻塞（因为自旋，并不挂起），提高了程序的响应速度。

    缺点：如果始终得不到锁，线程会一直卡着，消耗CPU（一般自旋就是 while true）

    适用：追求响应时间，同步块执行速度非常快

3. 重量级锁：严格排他的锁，也就是悲观锁。`synchronized` 就是这类。

    优点：线程竞争不使用自旋，不会消耗 CPU

    缺点：线程阻塞，响应时间慢

    适用：追求吞吐量，同步块执行时间较长

## happens-before 原则（先行发生原则）

并发编程需要考虑三个因素，分别是：

1. 原子性

    Java 虚拟机模型只保证**基本数据类型的读取和赋值**是原子的。因此例如 `int x = 1` 是一个原子操作。但例如 `int y = x` 或者 `x++` 这些既包含读取又包含写入（还包含计算的）就不是原子操作。

    使用 `synchronized` 或者 Lock 可以保证一系列操作的原子性，上面已经提过了。

2. 可见性

    当多个线程都需要使用某个变量的情况下，如果一个线程修改了这个变量，其他线程需要马上能读取到新的值，这称为可见性。

    在 Java 中可以使用 `volatile` 关键词来修饰一个变量。这样 Java 会保证修改的值会立即被更新到主存，当有其他线程需要读取时，它会去内存中读取新值。如果没有使用这个关键词修饰，那么它被修改之后，什么时候被写入主存是不确定的，当其他线程去读取时，此时内存中可能还是原来的旧值，因此无法保证可见性。

    如果使用 `synchronized` 或者 Lock 保证了只有一个线程能够操作，也能够保证可见性，此时就不需要使用 `volatile` 了。

3. 有序性

    Java 虚拟机允许对代码进行重排，保证最终执行效果和重排前一致。他主要判断代码之间的依赖关系来进行推断，因此在单线程情况下是没问题的。但是多线程时，重排可能会导致最终结果不同。

    同样使用 `volatile` 关键词可以保证一段代码的有序性，详细内容后面提到。

    如果使用 `synchronized` 或者 Lock 保证了只有一个线程能够操作，也能够保证有序性，此时就不需要使用 `volatile` 了。

    另外，Java 内存模型具备一些先天的“有序性”，即不需要通过任何手段就能够得到保证的有序性，这个通常也称为 happens-before 原则。如果两个操作的执行次序无法从 happens-before 原则推导出来，那么它们就不能保证它们的有序性，虚拟机会随意地对它们进行重排序。

happens-before 原则有 8 条：

1. 程序次序规则：一个线程内，按照代码顺序，书写在前面的操作先行发生于书写在后面的操作
2. 锁定规则：一个 unLock 操作先行发生于后面对同一个锁的 lock 操作

    也就是说要先释放锁才能下一次获得锁

3. volatile 变量规则：对一个变量的写操作先行发生于后面对这个变量的读操作

    `volatile` 的作用所在，保证一个变量先被更新写入，再被读取，才能读到新的值。

4. 传递规则：如果操作A先行发生于操作B，而操作B又先行发生于操作C，则可以得出操作A先行发生于操作C
5. 线程启动规则：Thread 对象的 start() 方法先行发生于此线程的每个一个动作
6. 线程中断规则：对线程 interrupt() 方法的调用先行发生于被中断线程的代码检测到中断事件的发生
7. 线程终结规则：线程中所有的操作都先行发生于线程的终止检测，我们可以通过 Thread.join() 方法结束、Thread.isAlive() 的返回值手段检测到线程已经终止执行
8. 对象终结规则：一个对象的初始化完成先行发生于他的 finalize() 方法的开始

## volatile

### 可见性

```java
//线程1
volatile boolean stop = false;
while (!stop) {
    doSomething();
}

//线程2
stop = true;
```

使用 `volatile` 修饰的 `stop` 变量后，当线程2对它进行修改时，会直接更新主存中的 `stop`，并通知所有线程内存中有 `stop` 的线程（例子中的线程1），这个 `stop` 已经被更新，你们的缓存已经失效。因此线程1在读取 `stop` 时，会重新从主存中读取，从而能够起到预期的效果（停止线程1的主体运行）。而如果不使用 `volatile`，在线程2更新 `stop` 时只更新线程内存。如果它一直没有更新主存中的 `stop`，那么线程1也感受不到这个变化，这就是可见性没有保证。

### 原子性

```java
public class Test {
    public volatile int inc = 0;

    public void increase() {
        inc++;
    }

    public static void main(String[] args) {
        final Test test = new Test();
        for(int i = 0;i < 10;i++){
            new Thread(){
                public void run() {
                    for(int j = 0;j < 1000;j++)
                        test.increase();
                };
            }.start();
        }

        while(Thread.activeCount()>1)  //保证前面的线程都执行完
            Thread.yield();
        System.out.println(test.inc);
    }
}
```

在运行这段代码之后，会发现预期效果是打印 10000，可实际上总是小于 10000。原因其实和线程同步是类似的。假设一种情况：

1. 线程1读取了 `inc`，假设值为10
2. 此时线程2插入，也读取 `inc`，也得到 10
3. 线程2累加，并写入，现在 `inc` 值为 11，并通知其他线程，这个 `inc` 的缓存失效
4. **重点在这里**线程1已经完成了读取 `inc`，所以虽然它被通知缓存失效，但因为它不需要再读取，所以对它已经不影响了。它接着操作累加，也只能得到 11

所以这里的 `volatile` 并不能保证原子性，和不写效果是一样的。要取得正确结果，需要使用 `synchroized` 或者 Lock 或者 `AtomaticInteger`，而不是 `volatile`。

### 有序性

被 `volatile` 修饰的变量的读写操作可以认为是一个分隔符，Java 重排指令时不能跨越分隔符，只能在分隔符两侧分别进行。举例来说：

```java
int x, y;
volatile flag;

x = 2; // 语句1
y = 0; // 语句2
flag = true; // 语句3
x = 10; // 语句4
y = 1; // 语句5
```

语句3 涉及了 `volatile` 变量的读写，因此它是一个分隔符。Java 重排指令时，不能把语句1，2换到3的后面，也不能把4，5换到3的前面。只能对换1和2，4和5的顺序，相当于语句3把程序隔开了。这样语句3保证了**执行到这一句的时候，语句1和2已经执行完成了，而语句4，5还没有执行**，所以语句1，2对于语句3，4，5都是可见的。

`volatile` 的底层实现是“lock前缀指令”，也叫内存屏障（或内存栅栏）。

### 适用场景

从三个性质的描述中可以看到，`synchronized` 和 Lock 能够保证三个性质，但会影响程序的执行效率。与之相比 `volatile` 的性能在某些情况下优于 `synchronized`，但两者并不是替代关系。使用 `volatile` 的场景必须具备以下 2 个条件：

1. 对变量的写操作不依赖于当前值

2. 该变量没有包含在具有其他变量的不变式中

这表明 `volatile` 变量独立于任何程序的状态，包括变量的当前状态，因此多用做状态控制量（flag）。

## 单例模式及并发版本

最普通的单例模式是这样的

```java
public class Singleton {
    private static Singleton instance;
    private Singleton (){}

    public static Singleton getInstance() {
     if (instance == null) {
         instance = new Singleton();
     }
     return instance;
    }
}
```

这在单线程情况下是OK的，但是多线程时，可能一个线程进入 `if` 后被挂起，第二个线程又进入 `if`，从而创建了多个实例。

一个简单的修改是把 `getInstance()` 变成 `synchronized`。

```java
public class Singleton {
    private static Singleton instance;
    private Singleton (){}

    public static synchronized Singleton getInstance() {
     if (instance == null) {
         instance = new Singleton();
     }
     return instance;
    }
}
```

这种虽然可以保证只创建一个实例，但过于严格导致性能下降。原因是，只有进 `if` 这段需要同步，当实例已经创建完成后，后续获取实例其实并非要求同步，可以并发获取的。

既然使用 `synchronized` 过于严格，我们就想到了 `volatile`。

```java
// 双重检验锁
public class Singleton{
	private volatile static Singleton instance;
	public static Singleton getSingleton() {
	    if (instance == null) {
	        synchronized (Singleton.class) {
	            if (instance == null) {
	                instance = new Singleton();
	            }
	        }
	    }
	    return instance ;
	}
}
```

这里有 2 个注意点：

1. 代码进行了 2 次 `if` 的判断。这是因为，当同步块仅限于内层 `if` 了之后，依然可能有多个线程同时进入外部的 `if`。所以如果去掉了内层的 `if`，那么依然会创建多个实例。而如果去掉了外面的 `if`，那结果又跟前面 `synchronized` 等价，性能下降。

2. `instance = new Singleton()` 这个操作并非原子操作，它大概包括 3 个步骤：

    1. 给 `instance` 分配内存
    2. 调用 `Singleton` 的构造函数来初始化成员变量
    3. 将 `instance` 对象指向分配的内存空间（执行完这步 `instance` 就不是 `null` 了）

    这里 2 和 3 是可以调换次序的。如果是采用 1-3-2 的顺序执行，在 3 执行完后线程被挂起，第二个线程进入，在外层 `if` 中判断已经不是 `null`，就直接 `return`，于是就得到一个部分初始化的实例（没有执行第 2 步）。

    因此需要给 `instance` 增加修饰词 `volatile`。这样如果其他线程要进入 `if` 判断，就涉及了 `instance` 的读取操作。这样它必须等待前一个线程执行完初始化后才能读取，就避免了上述问题的出现。

## 参考文档

* [synchornized 基础](https://juejin.im/post/594a24defe88c2006aa01f1c)
* [Java 并发编程的 volatile](https://www.cnblogs.com/dolphin0520/p/3920373.html)