using System;
using System.Collections.Generic;

namespace TestAppPhd
{
    internal class Program
    {
        static void Main(string[] args)
        {
            string line = Console.ReadLine();
            var spl = line.Split(' ');
            var n = int.Parse(spl[0]);
            var m = int.Parse(spl[1]);
            var set = new Set(n);

            for (int i = 0; i < m; i++)
            {
                line = Console.ReadLine();
                spl = line.Split(' ');
                
                set.Join(int.Parse(spl[0]) - 1, int.Parse(spl[1]) - 1);
            }

            
            for (int i = 0; i < n; i++)
            {
                for (int j = i + 1; j < n; j++)
                {
                    if (!set.Check(i, j))
                    {
                        Console.WriteLine($"{i + 1} {j + 1}");
                        return;
                    }
                }
            }

            Console.WriteLine("0");
        }

        public class Set
        {
            private int[] w;
            private int[] h;

            public Set(int n)
            {
                w = new int[n];
                h = new int[n];

                for (int i = 0; i < n; i++)
                {
                    w[i] = i;
                    h[i] = 1;
                }
            }

            public void Join(int a, int b)
            {
                a = FindRoot(a);
                b = FindRoot(b);

                if (a == b) return;

                if (h[a] > h[b])
                {
                    h[b] += h[a];
                    w[b] = a;
                }
                else
                {
                    h[a] += h[b];
                    w[a] = b;
                }
            }

            public bool Check(int a, int b)
            {
                var ra = FindRoot(a);
                w[a] = ra;
                h[a] = 2;


                var rb = FindRoot(b);
                w[b] = rb;
                h[b] = 2;

                return rb == ra;
            }

            private int FindRoot(int a)
            {
                while (w[a] != a)
                {
                    a = w[a];
                }

                return a;
            }
        }
    }
}
